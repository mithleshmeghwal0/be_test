package persistsql

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"regexp"

	"github.com/go-pg/migrations/v8"
	pg "github.com/go-pg/pg/v10"
)

type PersistSQL struct {
	db *pg.DB
}

func New(db *pg.DB) *PersistSQL {
	return &PersistSQL{
		db: db,
	}
}

func (p *PersistSQL) Migrate(fs fs.FS) (int64, int64, error) {
	c := migrations.NewCollection().
		DisableSQLAutodiscover(true).
		SetTableName("migrations")

	hfs := http.FS(fs)
	f, err := hfs.Open(".")
	if err != nil {
		if os.IsNotExist(err) {
			return 0, 0, nil
		}

		return 0, 0, err
	}

	list, err := f.Readdir(-1)
	if err != nil {
		return 0, 0, err
	}
	if len(list) == 0 {
		return 0, 0, nil
	}

	if err := c.DiscoverSQLMigrationsFromFilesystem(hfs, "migrations"); err != nil {
		return 0, 0, err
	}

	if _, _, err := c.Run(p.db, "init"); err != nil {
		return 0, 0, fmt.Errorf("init migrations: %v", err)
	}

	_, err = c.Version(p.db)
	if err != nil {
		return 0, 0, fmt.Errorf("read db version: %v", err)
	}

	oldV, newV, err := c.Run(p.db, "up")
	if err != nil {
		return 0, 0, fmt.Errorf("migrate: %v", err)
	}

	return oldV, newV, nil
}

func (p *PersistSQL) CreateResource(ctx context.Context, resource interface{}) (interface{}, error) {
	if err := p.db.WithContext(ctx).RunInTransaction(ctx, func(tx *pg.Tx) error {
		if _, err := tx.Model(resource).Insert(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return resource, nil
}

func (p *PersistSQL) GetResource(ctx context.Context, resource interface{}) (interface{}, error) {
	if err := p.db.ModelContext(ctx, resource).WherePK().Select(); err != nil {
		return nil, err
	}

	return resource, nil
}

func (p *PersistSQL) UpdateResource(ctx context.Context, resource interface{}, fields []string) (interface{}, error) {
	if err := p.db.WithContext(ctx).RunInTransaction(ctx, func(tx *pg.Tx) error {
		query := tx.Model(resource).Returning("*").Column("update_time")
		for _, col := range fields {
			query.Column(col)
		}

		query.WherePK()

		if _, err := query.Update(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return resource, nil
}

func (p *PersistSQL) DeleteResource(ctx context.Context, resource interface{}) (interface{}, error) {
	if err := p.db.WithContext(ctx).RunInTransaction(ctx, func(tx *pg.Tx) error {
		query := tx.Model(resource).WherePK().Returning("*").WherePK()

		if _, err := query.Delete(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return resource, nil
}

const maxPageSize = 1000

func (p *PersistSQL) ListResources(ctx context.Context, resource interface{}, filter string, offset, pageSize int) (int, error) {
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	query := p.db.ModelContext(ctx, resource)

	if filter != "" {
		f, err := parseFilter(filter)
		if err != nil {
			return 0, err
		}

		switch f.Operator {
		case "=":
			query = query.Where(fmt.Sprintf("%s = ?", f.Field), f.Value)
		case "!=":
			query = query.Where(fmt.Sprintf("%s != ?", f.Field), f.Value)
		case "LIKE":
			query = query.Where(fmt.Sprintf("%s ILIKE ?", f.Field), fmt.Sprintf("%%%v%%", f.Value))
		default:
			return 0, fmt.Errorf("invalid operator: %v", f.Operator)
		}
	}
	query.Offset(offset)
	query.Limit(pageSize)
	count, err := query.Count()
	if err != nil {
		return 0, err
	}

	if err := query.Select(resource); err != nil {
		return 0, err
	}

	return count, nil
}

var filterRegx = regexp.MustCompile(`^(.*)(=|!=|LIKE)"(.*)"$`)

type Filter struct {
	Field    string
	Operator string
	Value    string
}

func parseFilter(filter string) (*Filter, error) {
	str := filterRegx.FindStringSubmatch(filter)
	if len(str) != 4 {
		return nil, fmt.Errorf("invalid filter: %v", str)
	}

	return &Filter{
		Field:    str[1],
		Operator: str[2],
		Value:    str[3],
	}, nil
}
