package product

import (
	"context"
	"database/sql"
	"fmt"
	"product-service/internal/db/generated"

	"github.com/jmoiron/sqlx"
)

// Repository provides access to product data via sqlc-generated queries
type Repository struct {
	q *generated.Queries
}

// NewRepository creates a new Repository with a connected database
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{q: generated.New(db.DB)}
}

// ListProducts retrieves all products from the database
func (r *Repository) ListProducts(ctx context.Context) ([]generated.Product, error) {
	products, err := r.q.ListProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not list products: %w", err)
	}
	if products == nil {
		products = []generated.Product{}
	}
	return products, nil
}

// CreateProduct creates a product in the database
func (r *Repository) CreateProduct(ctx context.Context, name, description string, price string, stock int32) (generated.Product, error) {
	createProductParams := generated.CreateProductParams{
		Name: name,
		Description: sql.NullString{
			String: description,
			Valid:  description != "",
		},
		Price: price,
		Stock: stock,
	}
	product, err := r.q.CreateProduct(ctx, createProductParams)
	if err != nil {
		return generated.Product{}, fmt.Errorf("could not create product: %w", err)
	}

	return product, nil
}

// GetProduct retrieves a product from the database
func (r *Repository) GetProduct(ctx context.Context, id int32) (generated.Product, error) {
	product, err := r.q.GetProduct(ctx, id)
	if err != nil {
		return generated.Product{}, fmt.Errorf("could not get product: %w", err)
	}
	return product, nil
}

// UpdateProduct updates a product in the database
func (r *Repository) UpdateProduct(ctx context.Context, id int32, name, description string, price string, stock int32) (generated.Product, error) {
	updateProductParams := generated.UpdateProductParams{
		ID:   id,
		Name: name,
		Description: sql.NullString{
			String: description,
			Valid:  description != "",
		},
		Price: price,
		Stock: stock,
	}
	product, err := r.q.UpdateProduct(ctx, updateProductParams)
	if err != nil {
		return generated.Product{}, fmt.Errorf("could not update product: %w", err)
	}
	return product, nil
}

// DeleteProduct deletes a product from the database
func (r *Repository) DeleteProduct(ctx context.Context, id int32) error {
	err := r.q.DeleteProduct(ctx, id)
	if err != nil {
		return fmt.Errorf("could not delete product: %w", err)
	}
	return nil
}
