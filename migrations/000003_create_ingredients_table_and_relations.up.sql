CREATE TABLE IF NOT EXISTS ingredients (
	id bigserial PRIMARY KEY,
	created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
	name varchar(100) NOT NULL UNIQUE,
	food_type varchar(100)
);

CREATE TABLE IF NOT EXISTS recipes_ingredients (
	recipe_id bigint NOT NULL REFERENCES recipes ON DELETE CASCADE,
	ingredient_id bigint NOT NULL REFERENCES ingredients,
	amount varchar(100) NOT NULL,
	unit varchar(100),
	PRIMARY KEY (recipe_id, ingredient_id)
);
