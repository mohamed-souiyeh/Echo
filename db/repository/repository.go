package db

import (
	"context"
	"database/sql"
	sqlc "echo/db/sqlc_generated"
	"fmt"
)

// TODO: create an interface to make testability easy
type UserRepository interface {
	GetAllUsers(ctx context.Context) ([]sqlc.User, error)                                // for testing and debuging
	CreateUser(ctx context.Context, username string, password string) (sqlc.User, error) // Create
	// GetUserByUsername (ctx, username) // Read
	// GetUserById(ctx, id) // Read
	// SearchUsersByUsername(ctx, username) // Read - implement it using the LIKE keyword in the query.
}

type SQLiteUserRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

var _ UserRepository = (*SQLiteUserRepository)(nil)

//? TODO: implement the repository functions to run the test of seeding the users and printing them

func (r *SQLiteUserRepository) CreateUser(ctx context.Context, username string, hashedPassword string) (sqlc.User, error) {

	user, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Username: username,
		Password: hashedPassword,
	})

	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to create the user %s: %w", username, err)
	}

	return user, nil
}

func (r *SQLiteUserRepository) GetAllUsers(ctx context.Context) ([]sqlc.User, error) {
	users, err := r.queries.GetAllUsers(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			return []sqlc.User{}, nil
		}
		return nil, fmt.Errorf("failed getting all users: %w", err)
	}
	return users, nil
}

//? TODO: implement the seeding logic: DONE
// TODO: implement the proper repositories userRepository, roomRepository and msgRepository
