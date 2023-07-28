package sql

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"example.com/be_test/internal/models"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

type resultChan struct {
	result interface{}
	err    error
}

func (p *Persist) CreateUser(ctx context.Context, resource *models.User) (*models.User, error) {
	genRes, err := p.sql.CreateResource(ctx, resource)
	if err != nil {
		return nil, fmt.Errorf("%v : %w", err, errors.New("internal error"))
	}
	return genRes.(*models.User), nil
}

func (p *Persist) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	resource, err := p.sql.GetResource(ctx, &models.User{
		Common: models.Common{
			ID: id,
		},
	})
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("%v : %w", err, errors.New("user not found"))
		}

		return nil, fmt.Errorf("%v : %w", err, errors.New("internal error"))
	}

	return resource.(*models.User), nil
}

func (p *Persist) UpdateUser(ctx context.Context, resource *models.User, fields []string) (*models.User, error) {
	genRes, err := p.sql.UpdateResource(ctx, resource, fields)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("%v : %w", err, errors.New("user not found"))
		}

		return nil, fmt.Errorf("%v : %w", err, errors.New("internal error"))
	}

	return genRes.(*models.User), nil
}

func (p *Persist) DeleteUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	genRes, err := p.sql.DeleteResource(ctx, &models.User{
		Common: models.Common{
			ID: id,
		},
	})
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("%v : %w", err, errors.New("user not found"))
		}

		return nil, fmt.Errorf("%v : %w", err, errors.New("internal error"))
	}

	return genRes.(*models.User), nil
}

// expected format of filter is
// name="foo"
// name!="foo"

type Filter struct {
	Field    string
	Operator string
	Value    interface{}
}

const (
	defaultPageSize = 10
)

func (p *Persist) ListUser(ctx context.Context, filter string, pageSize int, nextPageToken string) ([]*models.User, string, error) {
	pageNumber := 0
	if nextPageToken != "" {
		npt, err := decodeToken(nextPageToken)
		if err != nil {
			return nil, "", fmt.Errorf("%v : %w", err, errors.New("invalid next page token"))
		}

		pageSize = npt.PageSize
		pageNumber = npt.PageNumber + 1
		nextPageToken = ""
	}

	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	if pageNumber < 0 {
		pageNumber = 0
	}

	users := []*models.User{}
	count, err := p.sql.ListResources(ctx, &users, filter, pageNumber*pageSize, pageSize)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, "", nil
		}

		return nil, "", fmt.Errorf("%v : %w", err, errors.New("internal error"))
	}

	if count > pageNumber*pageSize {
		nextPageToken, _ = encodeToken(pageNumber+1, pageSize)
	}

	return users, nextPageToken, nil
}

type NextPageToken struct {
	PageSize   int
	PageNumber int
}

func encodeToken(pageNumber, pageSize int) (string, error) {
	npt := &NextPageToken{
		PageSize:   pageSize,
		PageNumber: pageNumber,
	}

	data, err := json.Marshal(npt)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

func decodeToken(token string) (*NextPageToken, error) {
	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	npt := &NextPageToken{}
	err = json.Unmarshal(decodedToken, npt)
	if err != nil {
		return nil, err
	}

	return npt, nil
}
