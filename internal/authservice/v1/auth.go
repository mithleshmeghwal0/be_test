package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (a *AuthService) CreateAuth(ctx *gin.Context) {

	actor := uuid.New().String()

	token, err := a.jwt.BuildAndSignJWTToken(actor)
	if err != nil {
		a.log.WithError(err).Error("jwt.BuildAndSignJWTToken()")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errors.Unwrap(err).Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"token_type":   "JWT",
		"expires_in":   24 * time.Hour.Seconds(),
	})
}
