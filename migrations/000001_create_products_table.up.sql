-- Filename: migrations/000001_create_products_table.up.sql

-- CREATE TABLE IF NOT EXISTS comments (
--     id bigserial PRIMARY KEY,
--     created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
--     content text NOT NULL,
--     author text NOT NULL,
--     version integer NOT NULL DEFAULT 1
-- );

CREATE TABLE IF NOT EXISTS products (
    product_id bigserial PRIMARY KEY,
    name text NOT NULL,
    description text NOT NULL,
    category text NOT NULL,
    image_url text NOT NULL,
    price text NOT NULL,
    average_rating DECIMAL(3, 2) DEFAULT 0.00, 
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);

CREATE TABLE reviews (
    review_id bigserial PRIMARY KEY,
    product_id INT REFERENCES products(product_id) ON DELETE CASCADE,
    author VARCHAR(255),
    rating FLOAT CHECK (rating BETWEEN 1 AND 5),
    review_text text NOT NULL,
    helpful_count INT DEFAULT 0,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);

CREATE OR REPLACE FUNCTION automatic_average_rating()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE products
    SET average_rating = (
        SELECT ROUND(CAST(AVG(rating) AS NUMERIC), 2)
        FROM reviews
        WHERE reviews.product_id = NEW.product_id
    )
    WHERE product_id = NEW.product_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_product_rating
AFTER INSERT OR UPDATE OR DELETE ON reviews
FOR EACH ROW
EXECUTE FUNCTION automatic_average_rating();
    
