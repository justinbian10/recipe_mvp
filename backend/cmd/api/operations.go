package main

import (
	"context"
	"time"

	"recipemvp.justinbian/internal/data"
)

func (app *application) addFullRecipe(recipe *RecipeResource) error {
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
			Amount: ingredient.Amount,
			Unit: ingredient.Unit,
		}
		app.models.Ingredients.Insert(tx, ingredientData)
		ingredient.ID = ingredientData.ID
		app.models.Ingredients.AddForRecipe(tx, recipeData.ID, ingredientData)
	}

	for _, step := range recipe.Steps {
		stepData := &data.StepData{
			Description: step.Description,
		}
		app.models.Steps.Insert(tx, stepData)
		step.ID = stepData.ID
		app.models.Steps.AddForRecipe(tx, recipeData.ID, stepData)
	}

	return tx.Commit()
}

func (app *application) getFullRecipe(id int64) (*RecipeResource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tx, err := app.models.Recipes.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	recipeData, err := app.models.Recipes.Get(tx, id)
	if err != nil {
		return nil, err
	}
	ingredientDatas, err := app.models.Ingredients.GetForRecipe(tx, id)
	if err != nil {
		return nil, err
	}
	stepDatas, err := app.models.Steps.GetForRecipe(tx, id)
	if err != nil {
		return nil, err
	}
	var (
		ingredients []*IngredientResource
		steps []*StepResource
	)
	for _, ingredientData := range ingredientDatas {
		ingredients = append(ingredients, ingredientDataToResource(ingredientData))
	}
	for _, stepData := range stepDatas {
		steps = append(steps, stepDataToResource(stepData))
	}

	recipe := &RecipeResource{
		ID: recipeData.ID,
		Title: recipeData.Title,
		Description: recipeData.Description,
		ImageURL: recipeData.ImageURL,
		Servings: recipeData.Servings,
		CooktimeMinutes: recipeData.CooktimeMinutes,
		Ingredients: ingredients,
		Steps: steps,
		Version: recipeData.Version,
	}

	return recipe, tx.Commit()
}

func (app *application) getAllRecipes(title string, servings int, cooktime int, filters data.Filters) ([]*RecipeResource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tx, err := app.models.Recipes.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	recipeDatas, err := app.models.Recipes.GetAll(tx, title, servings, cooktime, filters)
	if err != nil {
		return nil, err
	}

	var recipeIDs []int64

	for _, recipeData := range recipeDatas {
		recipeIDs = append(recipeIDs, recipeData.ID)
	}

	ingredientDatas, ingredientRecipeIDs, err := app.models.Ingredients.GetForMultipleRecipes(tx, recipeIDs)
	if err != nil {
		return nil, err
	}
	ingredientMap := make(map[int64][]*IngredientResource)
	for i, ingredientData := range ingredientDatas {
		recipeID := ingredientRecipeIDs[i] 
		ingredientMap[recipeID] = append(ingredientMap[recipeID], ingredientDataToResource(ingredientData))
	}

	stepDatas, stepRecipeIDs, err := app.models.Steps.GetForMultipleRecipes(tx, recipeIDs)
	if err != nil {
		return nil, err
	}
	stepMap := make(map[int64][]*StepResource)
	for i, stepData := range stepDatas {
		recipeID := stepRecipeIDs[i] 
		stepMap[recipeID] = append(stepMap[recipeID], stepDataToResource(stepData))
	}

	var recipeResources []*RecipeResource
	for _, recipeData := range recipeDatas {
		recipe := &RecipeResource{
			ID: recipeData.ID,
			Title: recipeData.Title,
			Description: recipeData.Description,
			ImageURL: recipeData.ImageURL,
			Servings: recipeData.Servings,
			CooktimeMinutes: recipeData.CooktimeMinutes,
			Ingredients: ingredientMap[recipeData.ID],
			Steps: stepMap[recipeData.ID],
			Version: recipeData.Version,
		}
		recipeResources = append(recipeResources, recipe)
	}

	return recipeResources, nil
}

func ingredientDataToResource(ingredientData *data.IngredientData) *IngredientResource {
	return &IngredientResource{
		ID: ingredientData.ID,
		Name: ingredientData.Name,
		FoodType: ingredientData.FoodType,
		Amount: ingredientData.Amount,
		Unit: ingredientData.Unit,
	}
}

func ingredientResourceToData(ingredientResource *IngredientResource) *data.IngredientData {
	return &data.IngredientData{
		ID: ingredientResource.ID,
		Name: ingredientResource.Name,
		FoodType: ingredientResource.FoodType,
		Amount: ingredientResource.Amount,
		Unit: ingredientResource.Unit,
	}
}
		

func stepDataToResource(stepData *data.StepData) *StepResource {
	return &StepResource{
		ID: stepData.ID,
		Description: stepData.Description,
	}
}

func stepResourceToData(stepResource *StepResource) *data.StepData {
	return &data.StepData{
		ID: stepResource.ID,
		Description: stepResource.Description,
	}
}

func (app *application) updateFullRecipe(recipe *RecipeResource) error {
	recipeData := &data.RecipeData{
		ID: recipe.ID,
		Title: recipe.Title,
		Description: recipe.Description,
		ImageURL: recipe.ImageURL,
		Servings: recipe.Servings,
		CooktimeMinutes: recipe.CooktimeMinutes,
		Version: recipe.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tx, err := app.models.Recipes.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = app.models.Recipes.Update(tx, recipeData)
	recipe.Version = recipeData.Version
	if err != nil {
		return err
	}

	var (
		ingredients []*data.IngredientData
		steps []*data.StepData
	)

	if recipe.Ingredients != nil {
		for _, ingredient := range recipe.Ingredients {
			ingredients = append(ingredients, ingredientResourceToData(ingredient))
		}

		err = app.models.Ingredients.UpdateForRecipe(tx, recipeData.ID, ingredients)
		if err != nil {
			return err
		}

		index := 0	
		for _, ingredient := range ingredients {
			recipe.Ingredients[index].ID = ingredient.ID
			index += 1
		}
	}

	if recipe.Steps != nil {
		for _, step := range recipe.Steps {
			steps = append(steps, stepResourceToData(step))
		}

		err = app.models.Steps.UpdateForRecipe(tx, recipeData.ID, steps)
		if err != nil {
			return err
		}

		index := 0
		for _, step := range steps {
			recipe.Steps[index].ID = step.ID
			index += 1
		}
	}

	return tx.Commit()
}

func (app *application) deleteFullRecipe(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tx, err := app.models.Recipes.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = app.models.Recipes.Delete(tx, id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (app *application) getIngredients(name string) ([]string, error) {
	ingredientNames, err := app.models.Ingredients.GetAll(name)
	if err != nil {
		return nil, err
	}

	return ingredientNames, nil
}
