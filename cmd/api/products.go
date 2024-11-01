package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Duane-Arzu/test1/internal/data"
	_ "github.com/Duane-Arzu/test1/internal/data"
	"github.com/Duane-Arzu/test1/internal/validator"
	_ "github.com/Duane-Arzu/test1/internal/validator"
)

// Struct for handling incoming JSON for Product data
var incomingProductData struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Category    *string  `json:"category"`
	ImageURL    *string  `json:"image_url"`
	Price       *string  `json:"price"`
	AvgRating   *float32 `json:"avg_rating"`
}

func (a *applicationDependencies) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var incomingProductData struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Category    string `json:"category"`
		ImageURL    string `json:"image_url"`
		Price       string `json:"price"`
	}
	err := a.readJSON(w, r, &incomingProductData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	products := &data.Product{
		Name:        incomingProductData.Name,
		Description: incomingProductData.Description,
		Category:    incomingProductData.Category,
		ImageURL:    incomingProductData.ImageURL,
		Price:       incomingProductData.Price,
	}
	v := validator.New()
	data.ValidateProduct(v, products)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.productModel.InsertProduct(products)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("products/%d", products.ProductID))

	data := envelope{
		"Product": products,
	}
	err = a.writeJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) displayProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r, "pid")
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	products, err := a.productModel.GetProduct(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	data := envelope{
		"Product": products,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r, "pid")
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	products, err := a.productModel.GetProduct(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			a.notFoundResponse(w, r)
		} else {
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	var incomingProductData struct {
		Name        *string  `json:"name"`
		Description *string  `json:"description"`
		Category    *string  `json:"category"`
		ImageURL    *string  `json:"image_url"`
		Price       *string  `json:"price"`
		AvgRating   *float64 `json:"avg_rating"`
	}

	err = a.readJSON(w, r, &incomingProductData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if incomingProductData.Name != nil {
		products.Name = *incomingProductData.Name
	}
	if incomingProductData.Description != nil {
		products.Description = *incomingProductData.Description
	}
	if incomingProductData.Category != nil {
		products.Category = *incomingProductData.Category
	}
	if incomingProductData.ImageURL != nil {
		products.ImageURL = *incomingProductData.ImageURL
	}
	if incomingProductData.Price != nil {
		products.Price = *incomingProductData.Price
	}

	v := validator.New()
	data.ValidateProduct(v, products)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.productModel.UpdateProduct(products)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"Product": products,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r, "pid")
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.productModel.DeleteProduct(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.PIDnotFound(w, r, id)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	data := envelope{
		"message": "Product successfully deleted",
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) listProductHandler(w http.ResponseWriter, r *http.Request) {
	var queryParametersData struct {
		Name     string
		Category string
		data.Filters
	}

	queryParameters := r.URL.Query()
	queryParametersData.Name = a.getSingleQueryParameter(queryParameters, "name", "")
	queryParametersData.Category = a.getSingleQueryParameter(queryParameters, "category", "")

	v := validator.New()
	queryParametersData.Filters.Page = a.getSingleIntegerParameter(queryParameters, "page", 1, v)
	queryParametersData.Filters.PageSize = a.getSingleIntegerParameter(queryParameters, "page_size", 10, v)
	queryParametersData.Filters.Sort = a.getSingleQueryParameter(queryParameters, "sort", "product_id")
	queryParametersData.Filters.SortSafeList = []string{"product_id", "name", "-product_id", "-name"}

	data.ValidateFilters(v, queryParametersData.Filters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	products, metadata, err := a.productModel.GetAll(
		queryParametersData.Name,
		queryParametersData.Category,
		queryParametersData.Filters,
	)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	data := envelope{
		"products":  products,
		"@metadata": metadata,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
