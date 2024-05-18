package data

import (
	"context"
	"database/sql"
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

func (m StepModel) AddForRecipe(tx *sql.Tx, recipeID int64, stepID int64) error {
	query := `
		INSERT INTO recipes_steps VALUES ($1, $2)`
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		recipeID,
		stepID,
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

/*
func (m StepModel) InsertMany(steps []*StepData) error {
	query := `
		INSERT INTO steps (description)
			SELECT * FROM UNNEST($1)
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second())
	defer cancel()

	var (
		stepIDs []*int64
		stepDescriptions []string
	)
	for _, ingredient := range ingredients {
		stepIDs = append(stepIDs, *step.IDs)
		stepDescriptions = append(stepDescrption, step.Description)
	}

	return m.DB.QueryContext(ctx, query, stepDescriptions).Scan(stepIDs...)
}
*/
