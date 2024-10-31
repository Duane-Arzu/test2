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

	//setup routes for comments
	// router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthCheckHandler)
	// router.HandlerFunc(http.MethodPost, "/v1/comments", a.createCommentHandler)
	// router.HandlerFunc(http.MethodGet, "/v1/comments/:id", a.displayCommentHandler)
	// router.HandlerFunc(http.MethodPatch, "/v1/comments/:id", a.updateCommentHandler)
	// router.HandlerFunc(http.MethodDelete, "/v1/comments/:id", a.deleteCommentHandler)
	// router.HandlerFunc(http.MethodGet, "/v1/comments", a.listCommentHandler)
	// return a.recoverPanic(router)

//routes for Products
router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthCheckHandler)
router.HandlerFunc(http.MethodPost, "/v1/products", a.createProductHandler)
router.HandlerFunc(http.MethodGet, "/v1/products/:id", a.displayProductHandler)
router.HandlerFunc(http.MethodPatch, "/v1/products/:id", a.updateProductHandler)
router.HandlerFunc(http.MethodDelete, "/v1/products/:id", a.deleteProductHandler)
router.HandlerFunc(http.MethodGet, "/v1/products", a.listProductHandler)

//routes for Reviews
router.HandlerFunc(http.MethodPost, "/reviews", a.createReviewsHandler)
router.HandlerFunc(http.MethodGet, "/reviews/:rid", a.displayReviewsHandler)
router.HandlerFunc(http.MethodPatch, "/reviews/:rid", a.updateReviewsHandler)
router.HandlerFunc(http.MethodDelete, "/reviews/:rid", a.deleteReviewsHandler)
router.HandlerFunc(http.MethodGet, "/reviews", a.listReviewsHandler)

router.HandlerFunc(http.MethodGet, "/product-reviews/:rid", a.listProductReviewsHandler)
router.HandlerFunc(http.MethodGet, "/product/:id/reviews/:rid", a.getProductReviewsHandler)
router.HandlerFunc(http.MethodPatch, "/helpful-count/:rid", a.HelpfulCountHandler)




return a.recoverPanic(router)
}

