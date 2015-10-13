-- bactdb
-- Matthew R Dillon

CREATE TYPE e_roles AS ENUM('R', 'W', 'A');
-- 'R': read-only, default
-- 'W': read-write
-- 'A': administrator

CREATE TABLE users (
    id BIGSERIAL NOT NULL,
    email CHARACTER VARYING(254) NOT NULL UNIQUE,
    password CHARACTER(60) NOT NULL,
    name TEXT NOT NULL,
    role e_roles DEFAULT 'R' NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    CONSTRAINT users_pkey PRIMARY KEY (id)
);

