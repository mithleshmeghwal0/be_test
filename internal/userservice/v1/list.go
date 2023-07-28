package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (u *UserService) ListUsers(ctx *gin.Context) {
	filter := ctx.Query("filter")
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize"))
	pageToken := ctx.Query("nextPageToken")

	users, nextPageToken, err := u.persist.ListUser(ctx, filter, pageSize, pageToken)
	if err != nil {
		u.log.WithError(err).Error("ListUser()")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errors.Unwrap(err).Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users":         users,
		"nextPageToken": nextPageToken,
	})
}
