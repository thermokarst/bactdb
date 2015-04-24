-- bactdb
-- Matthew R Dillon

CREATE TABLE unit_types (
    id BIGSERIAL NOT NULL,
    name TEXT NOT NULL,
    symbol CHARACTER VARYING(10) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    CONSTRAINT unit_types_pkey PRIMARY KEY (id)
);

