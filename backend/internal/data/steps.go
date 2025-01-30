package data

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type StepData struct {
	ID int64 `json:"id"`
	CreatedAt time.Time `json:"-"`
	Description string `json:"description"`
}

type StepModel struct {
	DB *sql.DB
}

func (m StepModel) AddForRecipe(tx *sql.Tx, recipeID int64, step *StepData) error {
	query := `
		INSERT INTO recipes_steps VALUES ($1, $2)`
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		recipeID,
		step.ID,
	}

	_, err := tx.ExecContext(ctx, query, args...)
	return err
}


func (m StepModel) Insert(tx *sql.Tx, step *StepData) error {
	query := `
		WITH res AS(
			INSERT INTO steps (description) VALUES ($1)
			ON CONFLICT (description) DO NOTHING
			RETURNING id, created_at)
		SELECT * FROM res UNION SELECT id, created_at FROM steps WHERE description=$1`
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	return tx.QueryRowContext(ctx, query, step.Description).Scan(&step.ID, &step.CreatedAt)
}

func (m StepModel) GetForRecipe(tx *sql.Tx, recipeID int64) ([]*StepData, error) {
	query := `
		SELECT s.id, s.created_at, s.description
		FROM recipes_steps AS rs LEFT JOIN steps AS s ON rs.step_id = s.id
		WHERE rs.recipe_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := tx.QueryContext(ctx, query, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []*StepData
	for rows.Next() {
		var step StepData
		args := []any{
			&step.ID,
			&step.CreatedAt,
			&step.Description,
		}
		err := rows.Scan(args...)
		if err != nil {
			return nil, err
		}
		steps = append(steps, &step)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return steps, nil
}

func (m StepModel) GetForMultipleRecipes(tx *sql.Tx, recipeIDs []int64) ([]*StepData, []int64, error) {
	var idsAsString []string
	for _, id := range recipeIDs {
		idsAsString = append(idsAsString, strconv.Itoa(int(id)))
	}
	ids := strings.Join(idsAsString, ", ")

	query := fmt.Sprintf(`
		SELECT s.id, rs.recipe_id, s.created_at, s.description
		FROM recipes_steps AS rs LEFT JOIN steps AS s on rs.step_id = s.id
		WHERE rs.recipe_id IN (%s)`, ids)


	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var (
		steps []*StepData
		stepRecipeIDs []int64
	)
	for rows.Next() {
		var (
			step StepData
			recipeID int64
		)

		args := []any{
			&step.ID,
			&recipeID,	
			&step.CreatedAt,
			&step.Description,
		}
		err := rows.Scan(args...)
		if err != nil {
			return nil, nil, err
		}
		steps = append(steps, &step)
		stepRecipeIDs = append(stepRecipeIDs, recipeID)
	}
	
	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	return steps, stepRecipeIDs, nil
}


func (m StepModel) UpdateForRecipe(tx *sql.Tx, recipeID int64, steps []*StepData) error {
	err := m.DeleteForRecipe(tx, recipeID)
	if err != nil {
		return err
	}

	for _, step := range steps {
		err = m.Insert(tx, step)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO recipes_steps (recipe_id, step_id) VALUES `

	args := []any{recipeID}
	for i, step := range steps {
		query += fmt.Sprintf("($1, $%d),", i+2)
		args = append(args, step.ID) 
	}
	query = query[:len(query)-1]
	
	query += `
		ON CONFLICT(recipe_id, step_id) DO NOTHING`
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	_, err = tx.ExecContext(ctx, query, args...)
	return err
}

func (m StepModel) DeleteForRecipe(tx *sql.Tx, recipeID int64) error {
	query := `
		DELETE FROM steps WHERE steps.id IN
			(SELECT s.id FROM recipes_steps AS rs LEFT JOIN steps AS s ON rs.step_id = s.id
			WHERE rs.recipe_id = $1)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, recipeID)
	return err
}
