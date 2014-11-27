-- bactdb
-- Matthew R Dillon

CREATE TABLE unit_types (
    id BIGSERIAL NOT NULL,
    name CHARACTER VARYING(100),
    symbol CHARACTER VARYING(10),

    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT unit_types_pkey PRIMARY KEY (id)
);

