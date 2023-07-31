package sql_test

import (
	"context"
	"fmt"
	"testing"

	"example.com/be_test/internal/models"
	"example.com/be_test/internal/sql"
	"example.com/be_test/pkg/env"
	"example.com/be_test/pkg/logger"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	dbUrl = env.MustGet("TESTDB_URL")

	log          = logger.New()
	SeededUserID = uuid.MustParse("f3fa60c1-02a4-496a-8c9b-c5418c9d3e67")
	SeededUser   = models.User{
		Common: models.Common{
			ID: SeededUserID,
		},
		Name:     "a",
		Email:    "a@bc.c",
		CreateBy: "",
	}
)

func setupDB() *pg.DB {
	dbOpts, err := pg.ParseURL(dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	return pg.Connect(dbOpts)
}

func cleanAndCloseDB(db *pg.DB) {
	_, err := db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func initPersist(db *pg.DB) *sql.Persist {
	persist, err := sql.New(log, db)
	if err != nil {
		log.Fatal(err)
	}

	err = persist.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	return persist
}

func TestPersist_CreateUser(t *testing.T) {
	ctx := context.Background()

	db := setupDB()
	defer cleanAndCloseDB(db)

	persist := initPersist(db)

	type args struct {
		user *models.User
	}

	testUserID := uuid.New()

	tests := []struct {
		name      string
		args      args
		seeder    func(persist *sql.Persist)
		wantError error
		want      *models.User
	}{
		{
			name: "Success",
			seeder: func(persist *sql.Persist) {

			},
			args: args{
				user: &models.User{
					Common: models.Common{
						ID: testUserID,
					},
					Name:     "mith",
					Email:    "a@bc.c",
					CreateBy: "",
				},
			},
			wantError: nil,
			want: &models.User{
				Common: models.Common{
					ID:      testUserID,
					Version: 1,
				},
				Name:     "mith",
				Email:    "a@bc.c",
				CreateBy: "",
			},
		},
		{
			name: "reused ID error",
			seeder: func(persist *sql.Persist) {
				_, err := persist.CreateUser(ctx, &SeededUser)
				if err != nil {
					log.Panicf("persist.CreateUser(): %v", err)
				}
			},
			args: args{
				user: &models.User{
					Common: models.Common{
						ID: SeededUser.ID,
					},
					Name:     "mith",
					Email:    "a@bc.c",
					CreateBy: "",
				},
			},
			wantError: fmt.Errorf("ERROR #23505 duplicate key value violates unique constraint \"users_pkey\" : internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.seeder(persist)
			got, err := persist.CreateUser(ctx, tt.args.user)

			if tt.wantError != nil {
				if assert.Error(t, err) {
					assert.Equal(t, tt.wantError.Error(), err.Error())
					return
				}
				assert.Fail(t, "expected error")
				return
			}

			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Version, got.Version)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Email, got.Email)
		})
	}
}

func TestPersist_UpdateUser(t *testing.T) {
	ctx := context.Background()

	db := setupDB()
	defer cleanAndCloseDB(db)

	persist := initPersist(db)

	type args struct {
		user *models.User
	}

	tests := []struct {
		name      string
		args      args
		seeder    func(persist *sql.Persist)
		wantError error
		want      *models.User
	}{
		{
			name: "update existing user", // above test already created user in db, we will try to create it again
			seeder: func(persist *sql.Persist) {
				_, err := persist.CreateUser(ctx, &SeededUser)
				if err != nil {
					log.Panicf("persist.CreateUser(): %v", err)
				}
			},
			args: args{
				user: &models.User{
					Common: models.Common{
						ID:      SeededUser.ID,
						Version: 1,
					},
					Name:     "b",
					Email:    "b@bb.b",
					CreateBy: "",
				},
			},
			want: &models.User{
				Common: models.Common{
					ID:      SeededUser.ID,
					Version: 2,
				},
				Name:     "b",
				Email:    "b@bb.b",
				CreateBy: "",
			},
		},
		{
			name:   "update non existent user", // above test already created user in db, we will try to create it again
			seeder: func(persist *sql.Persist) {},
			args: args{
				user: &models.User{
					Common: models.Common{
						ID: uuid.New(),
					},
					Name:     "b",
					Email:    "b@bb.b",
					CreateBy: "",
				},
			},
			wantError: fmt.Errorf("pg: no rows in result set : user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.seeder(persist)
			got, err := persist.UpdateUser(ctx, tt.args.user, []string{"name", "email"})

			if tt.wantError != nil {
				if assert.Error(t, err) {
					assert.Equal(t, tt.wantError.Error(), err.Error())
					return
				}
				assert.Fail(t, "expected error")
				return
			}

			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Version, got.Version)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Email, got.Email)
		})
	}
}

func TestPersist_DeleteUser(t *testing.T) {
	ctx := context.Background()

	db := setupDB()
	defer cleanAndCloseDB(db)

	persist := initPersist(db)

	type args struct {
		userID uuid.UUID
	}

	tests := []struct {
		name      string
		args      args
		seeder    func(persist *sql.Persist)
		wantError error
	}{
		{
			name: "delete existing user",
			seeder: func(persist *sql.Persist) {
				_, err := persist.CreateUser(ctx, &SeededUser)
				if err != nil {
					log.Panicf("persist.CreateUser(): %v", err)
				}
			},
			args: args{
				userID: SeededUser.ID,
			},
			wantError: nil,
		},
		{
			name:   "delete non existent user",
			seeder: func(persist *sql.Persist) {},
			args: args{
				userID: uuid.New(),
			},
			wantError: fmt.Errorf("pg: no rows in result set : user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.seeder(persist)
			deleteUser, err := persist.DeleteUser(ctx, tt.args.userID)

			if tt.wantError != nil {
				if assert.Error(t, err) {
					assert.Equal(t, tt.wantError.Error(), err.Error())
					return
				}
				assert.Fail(t, "expected error")
				return
			}

			assert.NotNil(t, deleteUser.DeleteTime)
		})
	}
}

func TestPersist_GetUser(t *testing.T) {
	ctx := context.Background()

	db := setupDB()
	defer cleanAndCloseDB(db)

	persist := initPersist(db)

	type args struct {
		userID uuid.UUID
	}

	tests := []struct {
		name      string
		args      args
		seeder    func(persist *sql.Persist)
		wantError error
		want      *models.User
	}{
		{
			name: "get existing user",
			seeder: func(persist *sql.Persist) {
				_, err := persist.CreateUser(ctx, &SeededUser)
				if err != nil {
					log.Panicf("persist.CreateUser(): %v", err)
				}
			},
			args: args{
				userID: SeededUser.ID,
			},
			want: &models.User{
				Common: models.Common{
					ID:      SeededUser.ID,
					Version: 1,
				},
				Name:     "a",
				Email:    "a@bc.c",
				CreateBy: "",
			},
		},
		{
			name:   "get non existent user",
			seeder: func(persist *sql.Persist) {},
			args: args{
				userID: uuid.New(),
			},
			wantError: fmt.Errorf("pg: no rows in result set : user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.seeder(persist)
			got, err := persist.GetUser(ctx, tt.args.userID)

			if tt.wantError != nil {
				if assert.Error(t, err) {
					assert.Equal(t, tt.wantError.Error(), err.Error())
					return
				}
				assert.Fail(t, "expected error")
				return
			}

			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Version, got.Version)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Email, got.Email)
		})
	}
}

func TestPersist_ListUser(t *testing.T) {
	ctx := context.Background()

	db := setupDB()
	defer cleanAndCloseDB(db)

	persist := initPersist(db)

	// seeding
	_, err := persist.CreateUser(ctx, &SeededUser)
	if err != nil {
		log.Panicf("persist.CreateUser(): %v", err)
	}

	type args struct {
		filter string
	}

	tests := []struct {
		name      string
		args      args
		wantError error
		want      []*models.User
	}{
		{
			name: "list users",
			args: args{
				filter: "",
			},
			want: []*models.User{
				{
					Common: models.Common{
						ID:      SeededUser.ID,
						Version: 1,
					},
					Name:     "a",
					Email:    "a@bc.c",
					CreateBy: "",
				},
			},
		},
		{
			name: "filter list users",
			args: args{
				filter: `nameLIKE"a"`,
			},
			want: []*models.User{
				{
					Common: models.Common{
						ID:      SeededUser.ID,
						Version: 1,
					},
					Name:     "a",
					Email:    "a@bc.c",
					CreateBy: "",
				},
			},
		},
		{
			name: "filtered list users empty",
			args: args{
				filter: `nameLIKE"b"`,
			},
			want: []*models.User{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := persist.ListUser(ctx, tt.args.filter, 0, "")

			if tt.wantError != nil {
				if assert.Error(t, err) {
					assert.Equal(t, tt.wantError.Error(), err.Error())
					return
				}
				assert.Fail(t, "expected error")
				return
			}

			assert.Equal(t, len(tt.want), len(got))
		})
	}
}
