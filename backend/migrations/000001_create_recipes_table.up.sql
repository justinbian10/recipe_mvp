CREATE TABLE IF NOT EXISTS recipes (
	id bigserial PRIMARY KEY,
	created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
	title varchar(100) NOT NULL,
	description text NOT NULL,
	image_url text,
	servings integer NOT NULL,
	cooktime_minutes integer NOT NULL,
	version integer NOT NULL DEFAULT 1
);
