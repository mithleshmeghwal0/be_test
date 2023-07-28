package v1

import (
	"errors"
	"net/http"

	"example.com/be_test/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (u *UserService) CreateUser(ctx *gin.Context) {
	var resource models.User

	// Bind the JSON request body to the User struct
	if err := ctx.ShouldBindJSON(&resource); err != nil {
		u.log.WithError(err).Error("ShouldBindJSON()")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return
	}

	resource.ID = uuid.New()

	if err := resource.Validate(); err != nil {
		u.log.WithError(err).Error("resource.Validate()")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resource.ID = uuid.New()

	genRes, err := u.persist.CreateUser(ctx, &resource)
	if err != nil {
		u.log.WithError(err).Error("CreateUser()")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errors.Unwrap(err).Error()})
		return
	}

	ctx.JSON(http.StatusCreated, genRes)
}
