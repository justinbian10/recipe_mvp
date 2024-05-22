package main

import (
	"fmt"
	"errors"
	"net/http"

	"recipemvp.justinbian/internal/data"
)

type RecipeResource struct {
	ID int64 `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	ImageURL *string `json:"image_url,omitempty"`
	Ingredients []*IngredientResource `json:"ingredients"`
	Steps []*StepResource `json:"steps"`
	Servings int32 `json:"servings"`
	CooktimeMinutes int32 `json:"cooktime"`
	Version int32 `json:"version"`
}

type IngredientResource struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	FoodType string `json:"food_type,omitempty"`
	Amount string `json:"amount"`
	Unit string `json:"unit"`
}

type StepResource struct {
	ID int64 `json:"id"`
	Description string `json:"description"`
}

func (app *application) getRecipeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	recipe, err := app.getFullRecipe(id)
	
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"recipe": recipe}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createRecipeHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"Title"`
		Description string `json:"description"`
		Servings int32 `json:"servings"`
		CooktimeMinutes int32 `json:"cooktime"`
		Ingredients []*IngredientResource `json:"ingredients"`
		Steps []*StepResource `json:"steps"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	recipe := &RecipeResource{
		Title: input.Title,
		Description: input.Description,
		Servings: input.Servings,
		CooktimeMinutes: input.CooktimeMinutes,
		Ingredients: input.Ingredients,
		Steps: input.Steps,
	}

	err = app.addFullRecipe(recipe)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/recipes/%d", recipe.ID))
	
	err = app.writeJSON(w, http.StatusCreated, envelope{"recipe": recipe}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	recipe, err := app.getFullRecipe(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title *string `json:"title"`
		Description *string `json:"description"`
		ImageURL *string `json:"image_url"`
		Servings *int32 `json:"servings"`
		CooktimeMinutes *int32 `json:"cooktime"`
		Ingredients []*IngredientResource `json:"ingredients"`
		Steps []*StepResource `json:"steps"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Description != nil {
		recipe.Description = *input.Description
	}

	if input.ImageURL != nil {
		recipe.ImageURL = input.ImageURL
	}

	if input.Servings != nil {
		recipe.Servings = *input.Servings
	}

	if input.CooktimeMinutes != nil {
		recipe.CooktimeMinutes = *input.CooktimeMinutes
	}

	currIngredients := recipe.Ingredients 
	currSteps := recipe.Steps

	recipe.Ingredients = input.Ingredients
	recipe.Steps = input.Steps

	err = app.updateFullRecipe(recipe)

	if input.Ingredients == nil {
		recipe.Ingredients = currIngredients
	}
	
	if input.Steps == nil {
		recipe.Steps = currSteps
	}

	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"recipe": recipe}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteRecipeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.deleteFullRecipe(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "recipe deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
