package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *applicationDependences) routes() http.Handler {
	//setup a new router
	router := httprouter.New()

	//handle 405
	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)

	//method 404
	router.NotFound = http.HandlerFunc(a.notFoundResponse)

	//setup routes
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthCheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/comments", a.createCommentHandler)
	router.HandlerFunc(http.MethodGet, "/v1/comments/:id", a.displayCommentHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/comments/:id", a.updateCommentHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/comments/:id", a.deleteCommentHandler)
	router.HandlerFunc(http.MethodGet, "/v1/comments", a.listCommentHandler)
	// return a.recoverPanic(router)
//New
router.HandlerFunc(http.MethodPost, "/v1/products", a.createProductHandler)
router.HandlerFunc(http.MethodGet, "/v1/products/:id", a.displayProductHandler)
router.HandlerFunc(http.MethodPatch, "/v1/products/:id", a.updateProductHandler)
router.HandlerFunc(http.MethodDelete, "/v1/products/:id", a.deleteProductHandler)
router.HandlerFunc(http.MethodGet, "/v1/products", a.listProductHandler)

router.HandlerFunc(http.MethodPost, "/v1/products/:id/reviews", a.createReviewHandler)
router.HandlerFunc(http.MethodGet, "/v1/products/:id/reviews/:reviewId", a.displayReviewHandler)
router.HandlerFunc(http.MethodPatch, "/v1/products/:id/reviews/:reviewId", a.updateReviewHandler)
router.HandlerFunc(http.MethodDelete, "/v1/products/:id/reviews/:reviewId", a.deleteReviewHandler)
router.HandlerFunc(http.MethodGet, "/v1/products/:id/reviews", a.listReviewHandler)
return a.recoverPanic(router)
}

