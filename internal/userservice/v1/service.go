package v1

import (
	"example.com/be_test/internal/persist"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	log     *logrus.Entry
	persist persist.Persist
}

func New(log *logrus.Entry, persist persist.Persist) *UserService {
	return &UserService{
		log:     log,
		persist: persist,
	}
}
