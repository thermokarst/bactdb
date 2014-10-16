-- bactdb
-- Matthew R Dillon

CREATE TABLE strainsobsmeasurements (
    id BIGSERIAL NOT NULL,
    strainsobservations_id BIGINT,
    measurement_table CHARACTER VARYING(15),
    measurement_id BIGINT,

    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT strainsobsmeasurements_pkey PRIMARY KEY (id),
    FOREIGN KEY (strainsobservations_id) REFERENCES strainsobservations(id)
);

