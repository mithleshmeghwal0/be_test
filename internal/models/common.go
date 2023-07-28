package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Common struct {
	tableName struct{} `pg:",discard_unknown_columns"`

	ID         uuid.UUID  `json:"id" pg:",pk,type:uuid"`
	CreateTime time.Time  `json:"create_time" pg:",notnull"`
	UpdateTime time.Time  `json:"update_time" pg:",notnull"`
	DeleteTime *time.Time `json:"delete_time" pg:",soft_delete"`
	Version    uint64     `json:"-" pg:",notnull,default:1"`
}

func (m *Common) BeforeInsert(ctx context.Context) (context.Context, error) {
	if m.CreateTime.IsZero() {
		m.CreateTime = time.Now()
		m.UpdateTime = m.CreateTime
	}
	return ctx, nil
}

func (m *Common) BeforeUpdate(ctx context.Context) (context.Context, error) {
	m.UpdateTime = time.Now()
	return ctx, nil
}

func (m Common) ResourceVersion() uint64 {
	return m.Version
}

func (m Common) IsFieldOutputOnly(field string) bool {
	list := [...]string{
		"id",
		"create_time",
		"update_time",
		"delete_time",
		"version",
	}

	for _, curr := range list {
		if curr == field {
			return true
		}
	}

	return false
}
