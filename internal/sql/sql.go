package sql

import (
	"embed"
	"fmt"

	"example.com/be_test/pkg/persistsql"
	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Persist struct {
	db  *pg.DB
	sql *persistsql.PersistSQL
	log *logrus.Entry
}

func New(log *logrus.Entry, db *pg.DB) (*Persist, error) {
	persistsql := persistsql.New(db)
	oldV, newV, err := persistsql.Migrate(embedMigrations)
	if err != nil {
		log.WithError(err).Error("persistsql.Migrate(()")
		return nil, fmt.Errorf("persist.Migrate(): %v", err)
	}

	if oldV == newV {
		log.WithField("version", oldV).Debug("database up to date")
	} else {
		log.WithField("old_version", oldV).WithField("new_version", newV).Debugf("migrated version %d -> %d\n", oldV, newV)
	}

	return &Persist{
		db:  db,
		sql: persistsql,
		log: log,
	}, nil
}
