-- bactdb
-- Matthew R Dillon

CREATE TABLE observation_types (
    id BIGSERIAL NOT NULL,
    observation_type_name CHARACTER VARYING(100),

    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT observation_types_pkey PRIMARY KEY (id)
);

