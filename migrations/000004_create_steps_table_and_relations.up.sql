CREATE TABLE steps (
	id bigserial PRIMARY KEY,
	created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
	description text NOT NULL UNIQUE
);

CREATE TABLE recipes_steps (
	recipe_id bigint NOT NULL REFERENCES recipes ON DELETE CASCADE,
	step_id bigint NOT NULL REFERENCES steps ON DELETE CASCADE,
	PRIMARY KEY (recipe_id, step_id)
);
