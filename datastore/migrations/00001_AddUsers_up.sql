-- bactdb
-- Matthew R Dillon

CREATE TABLE users (
    id BIGSERIAL NOT NULL,
    username CHARACTER VARYING(100),

    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX username_idx
    ON users
    USING btree
    (username COLLATE pg_catalog."default");

