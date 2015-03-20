-- bactdb
-- Matthew R Dillon

CREATE TABLE text_measurement_types (
    id BIGSERIAL NOT NULL,
    text_measurement_name CHARACTER VARYING(100) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    CONSTRAINT text_measurements_pkey PRIMARY KEY (id)
);

