-- bactdb
-- Matthew R Dillon

CREATE TABLE measurements (
    id BIGSERIAL NOT NULL,
    strain_id BIGINT,
    observation_id BIGINT,
    text_measurement_type_id BIGINT NULL,
    measurement_value NUMERIC(6, 4) NULL,
    confidence_interval NUMERIC(6, 4) NULL,
    unit_type_id BIGINT NULL,

    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT strainsobsmeasurements_pkey PRIMARY KEY (id),
    FOREIGN KEY (strain_id) REFERENCES strains(id),
    FOREIGN KEY (observation_id) REFERENCES observations(id),
    FOREIGN KEY (text_measurement_type_id) REFERENCES text_measurement_types(id),
    FOREIGN KEY (unit_type_id) REFERENCES unit_types(id),
    CONSTRAINT exclusive_data_type CHECK (
        (text_measurement_type_id IS NOT NULL
            AND measurement_value IS NULL
            AND confidence_interval IS NULL
            AND unit_type_id IS NULL)
        OR
        (text_measurement_type_id IS NULL
            AND measurement_value IS NOT NULL))

);

