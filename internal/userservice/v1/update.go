package v1

import (
	"errors"
	"net/http"

	"example.com/be_test/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (u *UserService) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		u.log.WithError(err).Error("uuid.Parse()")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}

	var resource models.User
	// Bind the JSON request body to the User struct
	if err = ctx.ShouldBindJSON(&resource); err != nil {
		u.log.WithError(err).Error("ShouldBindJSON()")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return
	}
	resource.ID = userID
	if err := resource.Validate(); err != nil {
		u.log.WithError(err).Error("resource.Validate()")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedRes, err := u.persist.UpdateUser(ctx, &resource, []string{"name", "email"})
	if err != nil {
		u.log.WithError(err).Error("UpdateUser()")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errors.Unwrap(err).Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedRes)
}
