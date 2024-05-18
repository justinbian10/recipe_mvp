package main

import (
	"context"
	"time"

	"recipemvp.justinbian/internal/data"
)

func (app *application) addRecipe(recipe *RecipeResource) error {
	recipeData := &data.RecipeData{
		Title: recipe.Title,
		Description: recipe.Description,
		Servings: recipe.Servings,
		CooktimeMinutes: recipe.CooktimeMinutes,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tx, err := app.models.Recipes.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	app.models.Recipes.Insert(tx, recipeData)
	recipe.ID = recipeData.ID
	recipe.Version = recipeData.Version

	for _, ingredient := range recipe.Ingredients {
		ingredientData := &data.IngredientData{
			Name: ingredient.Name,
			FoodType: ingredient.FoodType,
		}
		app.models.Ingredients.Insert(tx, ingredientData)
		ingredient.ID = ingredientData.ID
		app.models.Ingredients.AddForRecipe(tx, recipeData.ID, ingredient.ID, ingredient.Amount, ingredient.Unit)
	}

	for _, step := range recipe.Steps {
		stepData := &data.StepData{
			Description: step.Description,
		}
		app.models.Steps.Insert(tx, stepData)
		step.ID = stepData.ID
		app.models.Steps.AddForRecipe(tx, recipeData.ID, step.ID)
	}

	return tx.Commit()
}
