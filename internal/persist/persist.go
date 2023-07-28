package persist

import (
	"context"

	"example.com/be_test/internal/models"
	"github.com/google/uuid"
)

type Persist interface {
	GetUser(ctx context.Context, id uuid.UUID) (*models.User, error)
	CreateUser(ctx context.Context, resource *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, resource *models.User, fields []string) (*models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) (*models.User, error)
	ListUser(ctx context.Context, filter string, pageSize int, nextPageToken string) ([]*models.User, string, error)
}
