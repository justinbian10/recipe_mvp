package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/recipes/:id", app.getRecipeHandler)
	router.HandlerFunc(http.MethodGet, "/v1/recipes", app.getAllRecipesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/recipes", app.createRecipeHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/recipes/:id", app.updateRecipeHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/recipes/:id", app.deleteRecipeHandler)

	router.HandlerFunc(http.MethodGet, "/v1/ingredients", app.getIngredientsHandler)

	return app.enableCORS(router)
}
