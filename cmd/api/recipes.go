package main

import (
	"net/http"
)

func (app *application) getRecipeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func (app *application) createRecipeHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) updateRecipeHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) deleteRecipeHandler(w http.ResponseWriter, r *http.Request) {

}
