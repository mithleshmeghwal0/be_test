package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (u *UserService) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		u.log.WithError(err).Error("uuid.Parse()")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}

	resource, err := u.persist.GetUser(ctx, userID)
	if err != nil {
		u.log.WithError(err).Error("GetUser()")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errors.Unwrap(err).Error()})
		return
	}

	ctx.JSON(http.StatusOK, resource)
}
