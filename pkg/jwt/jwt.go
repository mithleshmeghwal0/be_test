package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var expireTime = time.Now().Add(30 * 24 * time.Hour) // 24 hous

type Options struct {
	SignKey string
}

type JWT struct {
	opts Options
}

func New(opts Options) *JWT {
	return &JWT{
		opts: opts,
	}
}

type Claims struct {
	Data string
	jwt.RegisteredClaims
}

func (j *JWT) BuildAndSignJWTToken(data string) (string, error) {
	// Create claims with multiple fields populated
	claims := &Claims{
		data,
		jwt.RegisteredClaims{
			Issuer:    "test",
			Subject:   "somebody",
			ID:        uuid.New().String(),
			Audience:  []string{},
			ExpiresAt: jwt.NewNumericDate(expireTime),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.opts.SignKey))
	if err != nil {
		return "", fmt.Errorf("token.SignedString(): %v : %w", err, errors.New("internal error"))
	}

	return signedToken, nil
}

func (j *JWT) ParseAndValidateToken(tokenString string) (string, error) {
	c := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, c, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.opts.SignKey), nil
	})
	if err != nil {
		// if errors.Is(err, jwt.ErrTokenMalformed) {
		// } else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		// } else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		// } else {
		// }

		return "", fmt.Errorf("jwt.ParseWithClaims(): %v : %w", err, errors.New("invalid token"))
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	return c.Data, nil
}
