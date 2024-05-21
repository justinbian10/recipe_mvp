package data

import (
	"context"
	"database/sql"
	"fmt"
	"errors"
	"time"

	_ "github.com/lib/pq"
)

type RecipeData struct {
	ID int64 `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title string `json:"title"`
	Description string `json:"description"`
	ImageURL *string `json:"image_url,omitempty"`
	Servings int32 `json:"servings"`
	CooktimeMinutes int32 `json:"cooktime"`
	Version int32 `json:"version"`
}


type RecipeModel struct {
	DB *sql.DB
}

func (m RecipeModel) Insert(tx *sql.Tx, recipe *RecipeData) error {
	query := `
		INSERT INTO recipes (title, description, servings, cooktime_minutes)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []any{recipe.Title, recipe.Description, recipe.Servings, recipe.CooktimeMinutes}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return tx.QueryRowContext(ctx, query, args...).Scan(&recipe.ID, &recipe.CreatedAt, &recipe.Version)
}

func (m RecipeModel) Get(tx *sql.Tx, id int64) (*RecipeData, error) {
	if id < 0 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, title, description, image_url, servings, cooktime_minutes, version
		FROM recipes
		WHERE id = $1`

	var recipe RecipeData

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&recipe.ID,
		&recipe.CreatedAt,
		&recipe.Title,
		&recipe.Description,
		&recipe.ImageURL,
		&recipe.Servings,
		&recipe.CooktimeMinutes,
		&recipe.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &recipe, nil
}

func (m RecipeModel) Update(tx *sql.Tx, recipe *RecipeData) error {
	query := `
		UPDATE recipes
		SET title = $1, description = $2, image_url = $3, servings = $4, cooktime_minutes = $5, version = version + 1
		WHERE id = $6 AND version = $7
		RETURNING version`

	args := []any{
		recipe.Title,
		recipe.Description,
		recipe.ImageURL,
		recipe.Servings,
		recipe.CooktimeMinutes,
		recipe.ID,
		recipe.Version,
	}

	fmt.Println(recipe)
	fmt.Println(recipe.ImageURL)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, args...).Scan(&recipe.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m RecipeModel) Delete(tx *sql.Tx, id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM recipes
		WHERE id = $1`
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
