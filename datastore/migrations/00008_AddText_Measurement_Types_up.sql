-- bactdb
-- Matthew R Dillon

CREATE TABLE text_measurement_types (
    id BIGSERIAL NOT NULL,
    text_measurement_name CHARACTER VARYING(100),

    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT text_measurements_pkey PRIMARY KEY (id)
);

