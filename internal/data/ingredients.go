package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type IngredientData struct {
	ID int64 
	CreatedAt time.Time 
	Name string 
	FoodType string
	Amount string
	Unit string
}

type IngredientModel struct {
	DB *sql.DB
}

func (m IngredientModel) AddForRecipe(tx *sql.Tx, recipeID int64, ingredient *IngredientData) error {
	query := `
		INSERT INTO recipes_ingredients VALUES ($1, $2, $3, $4)`
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		recipeID,
		ingredient.ID,
		ingredient.Amount,
		ingredient.Unit,
	}

	_, err := tx.ExecContext(ctx, query, args...)
	return err
}
/*
func (m IngredientModel) AddManyForRecipe(tx *sql.Tx, recipeID int64, ingredients []*IngredientData) error {
	query := `
		INSERT INTO recipes_ingredients (recipe_id, ingredient_id, amount, unit) VALUES `

	values := []any{}
	values = append(values, recipeID)
	for i, ingredient := range ingredients {
		query += fmt.Sprintf("($1, %d, %d, %d),", i*4+1, i*4+2, i*4+3, i*4+4)
		values = append(values, ingredient.ID, ingredient.Amount, ingredient.Unit)
	}
	query = query[:len(query)-1]

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, values...)
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

func (m IngredientModel) UpdateForRecipe(tx *sql.Tx, recipeID int64, ingredients []*IngredientData) error {
	err := m.DeleteRelationshipForRecipe(tx, recipeID)
	if err != nil {
		return err
	}

	for _, ingredient := range ingredients {
		err = m.Insert(tx, ingredient)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO recipes_ingredients (recipe_id, ingredient_id, amount, unit) VALUES `

	args := []any{recipeID}
	for i, ingredient := range ingredients {
		query += fmt.Sprintf("($1, $%d, $%d, $%d),", i*3+2, i*3+3, i*3+4)
		args = append(args, ingredient.ID, ingredient.Amount, ingredient.Unit)
	}
	query = query[:len(query)-1]
/*	
	query += `
		ON CONFLICT(recipe_id, ingredient_id) DO UPDATE SET
		amount = EXCLUDED.amount, unit = EXCLUDED.unit`
		*/
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = tx.ExecContext(ctx, query, args...)
	return err
}



func (m IngredientModel) GetForRecipe(tx *sql.Tx, recipeID int64) ([]*IngredientData, error) {
	query := `
		SELECT i.id, i.created_at, i.name, i.food_type, ri.amount, ri.unit
		FROM recipes_ingredients AS ri LEFT JOIN ingredients AS i ON ri.ingredient_id = i.id
		WHERE ri.recipe_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := tx.QueryContext(ctx, query, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []*IngredientData
	for rows.Next() {
		var ingredient IngredientData
		args := []any{
			&ingredient.ID,
			&ingredient.CreatedAt,
			&ingredient.Name,
			&ingredient.FoodType,
			&ingredient.Amount,
			&ingredient.Unit,
		}
		err := rows.Scan(args...)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, &ingredient)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ingredients, nil
}

func (m IngredientModel) DeleteRelationshipForRecipe(tx *sql.Tx, recipeID int64) error {
	query := `
		DELETE FROM recipes_ingredients WHERE recipe_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, recipeID)
	return err
}
