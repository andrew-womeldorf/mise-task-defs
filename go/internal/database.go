package pkg

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/andrew-womeldorf/mise-task-defs/go/sql/gen/users"
	"github.com/google/uuid"
)

type Database struct {
	db      *sql.DB
	queries *users.Queries
}

func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	queries := users.New(db)

	return &Database{
		db:      db,
		queries: queries,
	}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) CreateUser(ctx context.Context, email, username string) error {
	return d.queries.CreateUser(ctx, users.CreateUserParams{
		ID:       uuid.New(),
		Email:    email,
		Username: username,
	})
}

func (d *Database) GetUser(ctx context.Context, id uuid.UUID) (users.User, error) {
	return d.queries.GetUser(ctx, id)
}

func (d *Database) GetUsers(ctx context.Context, limit int64) ([]users.User, error) {
	return d.queries.GetUsers(ctx, limit)
}
