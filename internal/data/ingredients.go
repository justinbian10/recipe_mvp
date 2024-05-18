package data

import (
	"context"
	"database/sql"
	"time"
)

type IngredientData struct {
	ID int64 `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name string `json:"name"`
	FoodType string `json:"food_type"`
}

type IngredientModel struct {
	DB *sql.DB
}

func (m IngredientModel) AddForRecipe(tx *sql.Tx, recipeID int64, ingredientID int64, amount string, unit string) error {
	query := `
		INSERT INTO recipes_ingredients VALUES ($1, $2, $3, $4)`
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		recipeID,
		ingredientID,
		amount,
		unit,
	}

	_, err := tx.ExecContext(ctx, query, args...)
	return err
}

/*
func (m IngredientModel) AddManyForRecipe(recipeID int64, ingredientIDs []int64, amounts []string, units []string) error {
	query := `
		INSERT INTO recipes_ingredients 
			SELECT $1, col1, col2, col3 FROM 
			UNNEST($2, $3, $4) AS u(col1, col2, col3)`
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		recipeID,
		ingredientIDs,
		amounts,
		units
	}

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}
*/

func (m IngredientModel) Insert(tx *sql.Tx, ingredient *IngredientData) error {
	query := `
		WITH res AS(INSERT INTO ingredients (name, food_type) VALUES ($1, $2)
			ON CONFLICT (name) DO NOTHING
			RETURNING id, created_at)
		SELECT * FROM res UNION SELECT id, created_at FROM ingredients WHERE name=$1`
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		ingredient.Name,
		ingredient.FoodType,
	}

	return tx.QueryRowContext(ctx, query, args...).Scan(&ingredient.ID, &ingredient.CreatedAt)
}

/*
func (m IngredientModel) InsertMany(ingredients []*IngredientData) error {
	query := `
		INSERT INTO ingredients (name, food_type)
			SELECT * FROM UNNEST($1, $2)
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		ingIDs []*int64
		ingNames, ingFoodTypes []string
	)
	for _, ingredient := range ingredients {
		ingIDs = append(ingIDs, *ingredient.IDs)
		ingNames = append(ingNames, ingredient.Name)
		ingFoodTypes = append(ingFoodTypes, ingredient.FoodType)
	}
	args := []any{
		ingNames,
		ingFoodTypes,
	}

	return m.DB.QueryContext(ctx, query, args...).Scan(ingIDs...)
}
*/
