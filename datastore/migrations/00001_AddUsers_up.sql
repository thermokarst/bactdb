-- bactdb
-- Matthew R Dillon

CREATE TABLE users (
    id BIGSERIAL NOT NULL,
    username CHARACTER VARYING(100),

    createdat TIMESTAMP WITH TIME ZONE,
    updatedat TIMESTAMP WITH TIME ZONE,
    deletedat TIMESTAMP WITH TIME ZONE,

    CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX username_idx
    ON users
    USING btree
    (username COLLATE pg_catalog."default");

