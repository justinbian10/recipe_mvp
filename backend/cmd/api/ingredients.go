package main

import (
	"net/http"
)

func (app *application) getIngredientsHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	ingredientNames, err := app.getIngredients(queryParams.Get("name"))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"ingredients": ingredientNames}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
