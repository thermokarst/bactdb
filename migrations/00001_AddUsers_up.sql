-- bactdb
-- Matthew R Dillon

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'e_roles') THEN
        CREATE TYPE e_roles AS ENUM('R', 'W', 'A');
        -- 'R': read-only, default
        -- 'W': read-write
        -- 'A': administrator
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL NOT NULL,
    email CHARACTER VARYING(254) NOT NULL UNIQUE,
    password CHARACTER(60) NOT NULL,
    name TEXT NOT NULL,
    role e_roles DEFAULT 'R' NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    CONSTRAINT users_pkey PRIMARY KEY (id)
);

