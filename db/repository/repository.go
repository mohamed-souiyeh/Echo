package db

import (
	"context"
	"database/sql"
	sqlc "echo/db/sqlc_generated"

	"github.com/charmbracelet/log"
)

// TODO: create an interface to make testability easy
type UserRepository interface {
	GetAllUsers(ctx context.Context) ([]sqlc.User, error) // for testing and debuging

	CreateUser(ctx context.Context, username string, hashedPassword string) (sqlc.User, error) // Create

	GetUserByUsername(ctx context.Context, username string) (sqlc.User, error) // Read
	// GetUserById(ctx, id) // Read
	// SearchUsersByUsername(ctx, username) // Read - implement it using the LIKE keyword in the query.
}

type PostgresUserRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

var _ UserRepository = (*PostgresUserRepository)(nil)

func (r *PostgresUserRepository) CreateUser(ctx context.Context, username string, hashedPassword string) (sqlc.User, error) {

	log.Debugf("repo db: %p, repo queries: %p", r.db, r.queries)

	user, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Username: username,
		Password: hashedPassword,
	})

	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (r *PostgresUserRepository) GetUserByUsername(ctx context.Context, username string) (sqlc.User, error) {
	user, err := r.queries.GetUserByUsername(ctx, username)

	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (r *PostgresUserRepository) GetAllUsers(ctx context.Context) ([]sqlc.User, error) {
	users, err := r.queries.GetAllUsers(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			return []sqlc.User{}, nil
		}
		return nil, err
	}
	return users, nil
}
