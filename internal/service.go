package internal

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"

	"example.com/be_test/internal/sql"
	"example.com/be_test/pkg/jwt"

	authSvc "example.com/be_test/internal/authservice/v1"
	userSvc "example.com/be_test/internal/userservice/v1"
)

type Service struct {
	srv  *http.Server
	db   *pg.DB
	log  *logrus.Entry
	opts Options
}

type Options struct {
	ListenAddressHTTP string
	Production        bool
	LogQuery          bool
	DBURL             string
	JWTSignKey        string
}

func New(ctx context.Context, opts Options, log *logrus.Entry) (*Service, error) {
	db, err := bootDB(opts, log)
	if err != nil {
		return nil, fmt.Errorf("bootDB(): %v", err)
	}

	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("db.Ping(): %v", err)
	}

	if opts.LogQuery {
		db.AddQueryHook(dbLogger{log})
	}

	persist, err := sql.New(log, db)
	if err != nil {
		return nil, fmt.Errorf("bootDB(): %v", err)
	}

	jwt := jwt.New(jwt.Options{
		SignKey: opts.JWTSignKey,
	})

	a := authSvc.New(log, jwt)

	h := userSvc.New(log, persist)

	engine := setupRoutes(log, opts, jwt, a, h)

	return &Service{
		srv: &http.Server{
			Addr:         opts.ListenAddressHTTP,
			Handler:      engine,
			ReadTimeout:  5 * time.Second,  // Maximum duration for reading the entire request
			WriteTimeout: 5 * time.Second,  // Maximum duration before timing out writes of the response
			IdleTimeout:  10 * time.Second, // Maximum amount of time to wait for the next request when keep-alives are enabled
		},
		db:   db,
		opts: opts,
		log:  log,
	}, nil
}

func (s *Service) Run() error {
	go func() {
		defer s.db.Close()
		s.log.Debugf("Listening and serving HTTP on %s\n", s.opts.ListenAddressHTTP)
		if err := s.srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				s.log.Debug(err)
				return
			}
			s.log.Error(err)
		}

	}()

	return nil
}

func (s *Service) Shutdown(ctx context.Context) {
	s.srv.Shutdown(ctx)
}

func setupRoutes(log *logrus.Entry, opts Options, jwt *jwt.JWT, a *authSvc.AuthService, u *userSvc.UserService) *gin.Engine {
	if opts.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "not found"})
	})
	engine.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"code": "not allowed"})
	})
	engine.Use(loggerMiddleware(log), gin.Recovery())

	authGroup := engine.Group("api/auth/v1")
	authGroup.POST("/auths", a.CreateAuth)

	usergroup := engine.Group("/api/user/v1")
	usergroup.Use(authenticate(jwt, log))
	usergroup.POST("/users", u.CreateUser)
	usergroup.GET("/users/:id", u.GetUser)
	usergroup.GET("/users", u.ListUsers)
	usergroup.PUT("/users/:id", u.UpdateUser)
	usergroup.DELETE("/users/:id", u.DeleteUser)

	return engine
}

func bootDB(opts Options, log *logrus.Entry) (*pg.DB, error) {
	dbOpts, err := pg.ParseURL(opts.DBURL)
	if err != nil {
		log.WithError(err).Error("pg.ParseURL()")
		return nil, fmt.Errorf("pg.ParseURL() :%v", err)
	}

	return pg.Connect(dbOpts), nil
}

type dbLogger struct {
	log *logrus.Entry
}

func (d dbLogger) BeforeQuery(ctx context.Context, q *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (d dbLogger) AfterQuery(ctx context.Context, q *pg.QueryEvent) error {
	bytes, err := q.FormattedQuery()
	if err == nil {
		d.log.Debug(string(bytes))
	}

	return nil
}

func authenticate(jwt *jwt.JWT, log *logrus.Entry) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		// Check if the Authorization header is empty or doesn't start with "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.WithField("authorization", authHeader).Debug("Invalid Authorization Header")
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		actor, err := jwt.ParseAndValidateToken(tokenString)
		if err != nil {
			log.WithField("tokenString", tokenString).WithError(err).Debug("jwt.ParseAndValidateToken()")
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		ctx.AddParam("actor", actor)

		ctx.Next()
	}
}

func loggerMiddleware(log *logrus.Entry) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()
		statusCode := c.Writer.Status()

		if statusCode < 200 && statusCode > 299 {
			latency := time.Since(startTime)
			clientIP := c.ClientIP()
			method := c.Request.Method
			path := c.Request.URL.Path

			log.WithFields(logrus.Fields{
				"status":    statusCode,
				"latency":   latency,
				"client_ip": clientIP,
				"method":    method,
				"path":      path,
			}).Debug("Request handled")
		}
	}
}
