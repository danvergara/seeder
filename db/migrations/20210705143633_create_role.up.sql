CREATE TABLE IF NOT EXISTS public.roles (
	id SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL UNIQUE
);
