ALTER TABLE recipes ADD CONSTRAINT recipes_servings_check CHECK (servings > 0);

ALTER TABLE recipes ADD CONSTRAINT recipes_cooktime_check CHECK (cooktime_minutes >= 0);
