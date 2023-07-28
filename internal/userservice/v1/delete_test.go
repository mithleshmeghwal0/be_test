package v1_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/be_test/internal/models"
	v1 "example.com/be_test/internal/userservice/v1"
	"example.com/be_test/pkg/logger"
	"example.com/be_test/tests/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteUser(t *testing.T) {
	var mockUUID = uuid.MustParse("b9500d3f-4fb8-4cbb-a530-4b144d9109a4")
	type args struct {
		id  string
		log *logrus.Entry
	}
	tests := []struct {
		name               string
		args               args
		mocks              func(df *depfields)
		expectedJSON       string
		expectedHttpStatus int
	}{
		{
			name: "success",
			args: args{
				id:  "b9500d3f-4fb8-4cbb-a530-4b144d9109a4",
				log: logger.New(),
			},
			mocks: func(df *depfields) {
				mockUser := &models.User{
					Name:     "test",
					Email:    "test@test.com",
					CreateBy: "",
				}
				mockUser.ID = mockUUID
				df.persistMock.On("DeleteUser", mock.Anything, mock.Anything).Return(mockUser, nil)
			},
			expectedHttpStatus: http.StatusOK,
			expectedJSON:       `{"id":"b9500d3f-4fb8-4cbb-a530-4b144d9109a4","create_time":"0001-01-01T00:00:00Z","update_time":"0001-01-01T00:00:00Z","delete_time":null,"name":"test","email":"test@test.com","create_by":""}`,
		},
		{
			name: "persist.DeleteUser() error",
			args: args{
				id:  "b9500d3f-4fb8-4cbb-a530-4b144d9109a4",
				log: logger.New(),
			},
			mocks: func(df *depfields) {
				df.persistMock.On("DeleteUser", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error: %w", errors.New("user not found")))
			},
			expectedHttpStatus: http.StatusInternalServerError,
			expectedJSON:       `{"error":"user not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := &depfields{
				persistMock: mocks.NewPersist(t),
			}
			tt.mocks(df)
			svc := v1.New(tt.args.log, df.persistMock)

			httpRecorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(httpRecorder)
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: mockUUID.String()})
			ctx.Request = httptest.NewRequest(http.MethodDelete, "/users/"+mockUUID.String(), nil)
			ctx.Request.Header.Set("Content-Type", "application/json")

			svc.DeleteUser(ctx)

			resp := httpRecorder.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedHttpStatus, resp.StatusCode)

			respBody, _ := io.ReadAll(resp.Body)
			assert.JSONEq(t, tt.expectedJSON, string(respBody))
		})
	}
}
