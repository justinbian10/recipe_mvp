package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict = errors.New("edit conflict")
)

type Models struct {
	Recipes RecipeModel
	Ingredients IngredientModel
	Steps StepModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Recipes: RecipeModel{DB: db},
		Ingredients: IngredientModel{DB: db},
		Steps: StepModel{DB: db},
	}
}
