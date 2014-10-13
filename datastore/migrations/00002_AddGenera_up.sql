-- bactdb
-- Matthew R Dillon

CREATE TABLE genera (
    id BIGSERIAL NOT NULL,
    genusname CHARACTER VARYING(100),

    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT genus_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX genusname_idx
    ON genera
    USING btree
    (genusname COLLATE pg_catalog."default");

