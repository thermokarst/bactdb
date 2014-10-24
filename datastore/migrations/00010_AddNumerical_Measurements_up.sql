-- bactdb
-- Matthew R Dillon

CREATE TABLE numerical_measurements (
    id BIGSERIAL NOT NULL,
    measurement_value NUMERIC(6, 4) NOT NULL,
    confidence_interval NUMERIC(6,4) NULL,
    unit_type_id BIGINT,

    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT numerical_measurements_pkey PRIMARY KEY (id),
    FOREIGN KEY (unit_type_id) REFERENCES unit_types(id)
);

