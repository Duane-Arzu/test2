package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Duane-Arzu/comments/internal/data"
	"github.com/Duane-Arzu/comments/internal/validator"
	_ "github.com/Duane-Arzu/comments/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type envelope map[string]any

func (a *applicationDependencies) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	jsResponse, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	jsResponse = append(jsResponse, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)
	_, err = w.Write(jsResponse)
	if err != nil {
		return err
	}

	return nil

}

func (a *applicationDependencies) readJSON(w http.ResponseWriter, r *http.Request,
	destination any) error {

	maxBytes := 256_000
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(destination)

	if err != nil {

		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("the body contains badly-formed JSON (at character %d)", syntaxError.Offset)
			// Decode can also send back an io error message
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("the body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("the body contains the incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("the body contains the incorrect  JSON type (at character %d)",
				unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("the body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(),
				"json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("the body must not be larger that %d bytes", maxBytesError.Limit)
		case errors.Is(err, io.EOF):
			return errors.New("the body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("the body must not be larger that %d bytes", maxBytesError.Limit)

		case errors.Is(err, io.EOF):
			return errors.New("the body must not be empty")

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})

	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (a *applicationDependencies) readIDParam(r *http.Request) (int64, error) {

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (a *applicationDependencies) getSingleQueryParameter(queryParameters url.Values, key string, defaultValue string) string {

	result := queryParameters.Get(key)
	if result == "" {
		return defaultValue
	}
	return result
}

func (a *applicationDependencies) getMultipleQueryParameters(queryParameters url.Values, key string, defaultValue []string) []string {

	result := queryParameters.Get(key)
	if result == "" {
		return defaultValue
	}
	return strings.Split(result, ",")
}

func (a *applicationDependencies) getSingleIntegerParameter(queryParameters url.Values, key string, defaultValue int, v *validator.Validator) int {

	result := queryParameters.Get(key)
	if result == "" {
		return defaultValue
	}
	// try to convert to an integer
	intValue, err := strconv.Atoi(result)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return intValue
}

// Write JSON response for a product
func (a *applicationDependences) writeProductResponse(w http.ResponseWriter, status int, product data.Product) error {
	return a.writeJSON(w, status, envelope{"product": product}, nil)
}

// Write JSON response for a list of products
func (a *applicationDependences) writeProductsResponse(w http.ResponseWriter, status int, products []data.Product) error {
	return a.writeJSON(w, status, envelope{"products": products}, nil)
}

// Write JSON response for a review
func (a *applicationDependences) writeReviewResponse(w http.ResponseWriter, status int, review data.Review) error {
	return a.writeJSON(w, status, envelope{"review": review}, nil)
}

// Write JSON response for a list of reviews
func (a *applicationDependences) writeReviewsResponse(w http.ResponseWriter, status int, reviews []data.Review) error {
	return a.writeJSON(w, status, envelope{"reviews": reviews}, nil)
}

func (a *applicationDependences) validateProduct(product data.Product) error {
	v := validator.New()
	if product.Name == "" {
		v.AddError("name", "Name is required.")
	}
	if product.Category == "" {
		v.AddError("category", "Category is required.")
	}
	if product.ImageURL == "" {
		v.AddError("image_url", "Image URL is required.")
	}
	if product.AverageRating < 0 || product.AverageRating > 5 {
		v.AddError("average_rating", "Average rating must be between 0 and 5.")
	}
	if !v.Valid() {
		return fmt.Errorf("validation error: %s", strings.Join(v.Errors, ", "))
	}
	return nil
}

func (a *applicationDependences) validateReview(review data.Review) error {
	v := validator.New()
	if review.Rating < 1 || review.Rating > 5 {
		v.AddError("rating", "Rating must be between 1 and 5.")
	}
	if review.ReviewText == "" {
		v.AddError("review_text", "Review text is required.")
	}
	if !v.Valid() {
		return fmt.Errorf("validation error: %s", strings.Join(v.Errors, ", "))
	}
	return nil
}

// Extract search, filter, sort parameters for products
func (a *applicationDependences) getProductQueryParams(query url.Values) (string, []string, string) {
	search := a.getSingleQueryParameter(query, "search", "")
	filters := a.getMultipleQueryParameters(query, "filter", []string{})
	sort := a.getSingleQueryParameter(query, "sort", "name")
	return search, filters, sort
}

// Extract search, filter, sort parameters for reviews
func (a *applicationDependences) getReviewQueryParams(query url.Values) (string, []string, string) {
	search := a.getSingleQueryParameter(query, "search", "")
	filters := a.getMultipleQueryParameters(query, "filter", []string{})
	sort := a.getSingleQueryParameter(query, "sort", "rating")
	return search, filters, sort
}

func (a *applicationDependences) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var product data.Product

	// Parse JSON input into product struct
	if err := a.readJSON(w, r, &product); err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Validate the product data
	if err := a.validateProduct(product); err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Set default average rating
	product.AverageRating = 0

	// Add the product to the database (example SQL function, assume it returns product with ID set)
	newProduct, err := a.commentModel.InsertProduct(product)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Send response with the newly created product
	if err := a.writeProductResponse(w, http.StatusCreated, newProduct); err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
