package user

import (
	"context"
	"fmt"
	"user-service/internal/db/generated"

	"github.com/jmoiron/sqlx"
)

// Repository provides access to user data via sqlc-generated queries
type Repository struct {
	q *generated.Queries
}

// NewRepository creates a new Repository with a connected database
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{q: generated.New(db.DB)}
}

// ListUsers retrieves all users from the database
func (r *Repository) ListUsers(ctx context.Context) ([]generated.User, error) {
	users, err := r.q.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not list users: %w", err)
	}
	if users == nil {
		users = []generated.User{}
	}
	return users, nil
}


// CreateUsers creates a user to the database
func (r *Repository) CreateUser(ctx context.Context, name, email string) (generated.User, error) {
	createUserParams := generated.CreateUserParams {
		Name: name,
		Email: email,
	}
	user, err := r.q.CreateUser(ctx, createUserParams)
	if err != nil {
		return generated.User{}, fmt.Errorf("could not create user: %w", err)
	}

	return user, nil
}

// GetUser retrieves a user from the database
func (r *Repository) GetUser(ctx context.Context, id int32) (generated.User, error) {
	user, err := r.q.GetUser(ctx, id)
	if err != nil {
		return generated.User{}, fmt.Errorf("could not get user: %w", err)
	}
	return user, nil
}

// UpdateUser updates a user in the database
func (r *Repository) UpdateUser(ctx context.Context, id int32, name, email string) (generated.User, error) {
	updateUserParams := generated.UpdateUserParams {
		ID: id,
		Name: name,
		Email: email,
	}
	user, err := r.q.UpdateUser(ctx, updateUserParams)
	if err != nil {
		return generated.User{}, fmt.Errorf("could not update user: %w", err)
	}
	return user, nil
}

// DeleteUser deletes a user from the database
func (r *Repository) DeleteUser(ctx context.Context, id int32) error {
	err := r.q.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("could not delete user: %w", err)
	}
	return nil
}

