package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/recipes/:recipeId", app.getRecipeHandler)
	router.HandlerFunc(http.MethodPost, "/v1/recipes/:recipeId", app.createRecipeHandler)
	router.HandlerFunc(http.MethodPut, "/v1/recipes/:recipeId", app.updateRecipeHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/recipes/:recipeId", app.deleteRecipeHandler)

	return router
}
