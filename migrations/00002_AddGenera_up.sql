-- bactdb
-- Matthew R Dillon

CREATE TABLE genera (
    id BIGSERIAL NOT NULL,
    genus_name TEXT NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    CONSTRAINT genus_pkey PRIMARY KEY (id)
);

