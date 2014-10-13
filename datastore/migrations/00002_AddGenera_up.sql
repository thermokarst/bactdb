-- bactdb
-- Matthew R Dillon

CREATE TABLE genera (
    id BIGSERIAL NOT NULL,
    genusname CHARACTER VARYING(100),

    createdat TIMESTAMP WITH TIME ZONE,
    updatedat TIMESTAMP WITH TIME ZONE,
    deletedat TIMESTAMP WITH TIME ZONE,

    CONSTRAINT genus_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX genusname_idx
    ON genera
    USING btree
    (genusname COLLATE pg_catalog."default");

