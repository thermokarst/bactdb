-- bactdb
-- Matthew R Dillon

CREATE TABLE strainsobsmeasurements (
    id BIGSERIAL NOT NULL,
    strainsobservations_id BIGINT,
    text_measurement_type_id BIGINT NULL,
    measurement_value NUMERIC(6, 4) NULL,
    confidence_interval NUMERIC(6, 4) NULL,
    unit_type_id BIGINT NULL,

    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT strainsobsmeasurements_pkey PRIMARY KEY (id),
    FOREIGN KEY (strainsobservations_id) REFERENCES strainsobservations(id),
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

