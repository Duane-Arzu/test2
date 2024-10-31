// internal/data/models.go
package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	ImageURL    string    `json:"image_url"`
	AvgRating   float64   `json:"avg_rating"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Review struct {
	ID           string    `json:"id"`
	ProductID    string    `json:"product_id"`
	UserID       string    `json:"user_id"`
	Rating       int       `json:"rating"`
	Comment      string    `json:"comment"`
	HelpfulCount int       `json:"helpful_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ProductModel struct {
	DB *sql.DB
}

type ReviewModel struct {
	DB *sql.DB
}

// Update the applicationDependences struct in main.go
type applicationDependences struct {
	config       serverConfig
	logger       *slog.Logger
	productModel ProductModel
	reviewModel  ReviewModel
}

// ProductModel methods
func (p ProductModel) Insert(product *Product) error {
	query := `
        INSERT INTO products (name, description, category, image_url)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at`

	args := []interface{}{product.Name, product.Description, product.Category, product.ImageURL}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
}

func (p ProductModel) Get(id string) (*Product, error) {
	query := `
        SELECT id, name, description, category, image_url, avg_rating, created_at, updated_at
        FROM products
        WHERE id = $1`

	var product Product

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := p.DB.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Category,
		&product.ImageURL,
		&product.AvgRating,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &product, nil
}

// ReviewModel methods
func (r ReviewModel) Insert(review *Review) error {
	query := `INSERT INTO reviews (product_id, user_id, rating, comment)
              VALUES ($1, $2, $3, $4)
              RETURNING id, created_at, updated_at`

	args := []interface{}{review.ProductID, review.UserID, review.Rating, review.Comment}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&review.ID, &review.CreatedAt, &review.UpdatedAt)
	if err != nil {
		return err
	}

	// Update average rating in products table
	_, err = r.DB.ExecContext(ctx, `
        UPDATE products 
        SET avg_rating = (SELECT AVG(rating) FROM reviews WHERE product_id = $1)
        WHERE id = $1`, review.ProductID)

	return err
}

// handlers.go
func (app *applicationDependences) createProduct(w http.ResponseWriter, r *http.Request) {
	var product Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.productModel.Insert(&product)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/products/%s", product.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"product": product}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// routes.go
func (app *applicationDependences) routes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/v1/products", app.createProduct).Methods("POST")
	// Add other routes...

	return router
}

func (p ProductModel) GetAll(name, category, sort string) ([]Product, error) {
	query := `SELECT id, name, description, category, image_url, avg_rating, created_at, updated_at
              FROM products
              WHERE ($1 = '' OR name ILIKE '%' || $1 || '%') AND ($2 = '' OR category = $2)
              ORDER BY CASE WHEN $3 = 'name' THEN name END ASC, 
                       CASE WHEN $3 = 'avg_rating' THEN avg_rating END DESC`
	args := []interface{}{name, category, sort}

	rows, err := p.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.ID, &product.Name, &product.Description, &product.Category, &product.ImageURL,
			&product.AvgRating, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

// Initialize DB Schema
const productsTable = `
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    category TEXT NOT NULL,
    image_url TEXT NOT NULL,
    avg_rating DECIMAL(3,2) DEFAULT 0,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT NOT NULL,
    helpful_count INTEGER DEFAULT 0,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);`
