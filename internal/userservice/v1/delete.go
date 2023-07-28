package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (u *UserService) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		u.log.WithError(err).Error("uuid.Parse()")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}

	deletedRes, err := u.persist.DeleteUser(ctx, userId)
	if err != nil {
		u.log.WithError(err).Error("UpdateUser()")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errors.Unwrap(err).Error()})
		return
	}

	ctx.JSON(http.StatusOK, deletedRes)
}
