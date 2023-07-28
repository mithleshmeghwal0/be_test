package v1

import (
	"example.com/be_test/pkg/jwt"
	"github.com/sirupsen/logrus"
)

type AuthService struct {
	log *logrus.Entry
	jwt *jwt.JWT
}

func New(log *logrus.Entry, jwt *jwt.JWT) *AuthService {
	return &AuthService{
		log: log,
		jwt: jwt,
	}
}
