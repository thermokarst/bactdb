-- bactdb
-- Matthew R Dillon

CREATE TABLE users (
    id BIGSERIAL NOT NULL,
    username CHARACTER VARYING(100) NOT NULL,
    password CHARACTER VARYING(100) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX username_idx
    ON users
    USING btree
    (username COLLATE pg_catalog."default");

