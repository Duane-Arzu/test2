// internal/data/products.go
package data

import (
    "context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Duane-Arzu/test1/internal/validator"
)

type Product struct {
    ProductID    int64     `json:"product_id"`
    Name         string    `json:"name"`
    Description  string    `json:"description"`
    Category     string    `json:"category"`
    ImageURL     string    `json:"image_url"`
    Price        string    `json:"price"`
    AvgRating    float64   `json:"avg_rating"`
    CreatedAt    time.Time `json:"created_at"`
    Version   int32     `json:"version"`
}

// type Review struct {
//     ReviewID     int64     `json:"review_id"`
//     ProductID    int64     `json:"product_id"`
//     Author       string    `json:"author"`
//     Rating       int64     `json:"rating"`
//     Comment      string    `json:"comment"`
//     HelpfulCount int       `json:"helpful_count"`
//     CreatedAt    time.Time `json:"created_at"`
//     Version      int32     `json:"version"`
// }

type ProductModel struct {
    DB *sql.DB
}

// type ReviewModel struct {
//     DB *sql.DB
// }

// Validation function for Product struct
func ValidateProduct(v *validator.Validator, product *Product) {
	v.Check(product.Name != "", "name", "this is required")
	v.Check(len(product.Name) <= 100, "name", "cannot be more than 100 characters long")
	v.Check(product.Description != "", "description", "this is required")
	v.Check(len(product.Description) <= 500, "description", "cannot be more than 500 characters long")
	v.Check(product.Category != "", "category", "this is required")
	v.Check(product.ImageURL != "", "image_url", "this is required")
	v.Check(len(product.ImageURL) <= 255, "image_url", "cannot be more than 255 characters long")
	v.Check(len(product.Price) <= 10, "price", "cannot be more than 10 characters long")
	v.Check(product.Description != "", "description", "this is required")
}

// // Insert Row to comments table
// // expects a pointer to the actual comment content
func (p ProductModel) InsertProduct(product *Product) error {
    //the sql query to be executed against the database table
	query := `
		INSERT INTO products (name, description, category, image_url, price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING product_id, created_at, version
	`
    //the actual values to be passed into $1, $2, $3, $4 and $5
	args := []any{product.Name, product.Description, product.Category, product.ImageURL, product.Price}


 	// Create a context with a 3-second timeout. No database
 	// operation should take more than 3 seconds or we will quit it
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(
		&product.ProductID,
		&product.CreatedAt,
		&product.Version,
	)
}

// get a product from DB based on ID
func (p ProductModel) GetProduct(id int64) (*Product, error) {
	//check if the id is valid
    if id < 1 {
		return nil, ErrRecordNotFound
	}

    //the sql query to be excecuted against the database table
	query := `
		SELECT product_id, name, description, category, image_url, price, average_rating, created_at, version
		FROM products
		WHERE product_id = $1
	`

	var product Product
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := p.DB.QueryRowContext(ctx, query, id).Scan(
		&product.ProductID,
		&product.Name,
		&product.Description,
		&product.Category,
		&product.ImageURL,
		&product.Price,
		&product.AverageRating,
		&product.CreatedAt,
		&product.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &product, nil
}

func (p ProductModel) UpdateProduct(product *Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, category = $3, image_url = $4, price = $5, average_rating = $6, version = version + 1
		WHERE product_id = $7
		RETURNING version
	`

	// Removed `product.UpdatedAt` from the args slice
	args := []any{product.Name, product.Description, product.Category, product.ImageURL, product.Price, product.AverageRating, product.ProductID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&product.Version)
}

func (p ProductModel) DeleteProduct(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM products
		WHERE product_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := p.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (p ProductModel) GetAll(name string, category string, filters Filters) ([]*Product, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), product_id, name, description, category, image_url, price, average_rating, created_at, version
		FROM products
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '') 
		AND (to_tsvector('simple', category) @@ plainto_tsquery('simple', $2) OR $2 = '') 
		ORDER BY %s %s, product_id ASC 
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := p.DB.QueryContext(ctx, query, name, category, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()
	totalRecords := 0
	products := []*Product{}

	for rows.Next() {
		var product Product
		err := rows.Scan(
			&totalRecords,
			&product.ProductID,
			&product.Name,
			&product.Description,
			&product.Category,
			&product.ImageURL,
			&product.Price,
			&product.AverageRating,
			&product.CreatedAt,
			&product.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		products = append(products, &product)
	}

	err = rows.Err()
	if err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	return products, metadata, nil
}