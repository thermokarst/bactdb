-- bactdb
-- Matthew R Dillon

CREATE TABLE test_methods (
    id BIGSERIAL NOT NULL,
    name CHARACTER VARYING(100) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    CONSTRAINT test_methods_pkey PRIMARY KEY (id)
);

