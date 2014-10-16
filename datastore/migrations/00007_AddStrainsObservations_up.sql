-- bactdb
-- Matthew R Dillon

CREATE TABLE strainsobservations (
    id BIGSERIAL NOT NULL,
    strain_id BIGINT NOT NULL,
    observations_id BIGINT NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT strainsobservations_pkey PRIMARY KEY (id),
    FOREIGN KEY (strain_id) REFERENCES strains(id),
    FOREIGN KEY (observations_id) REFERENCES observations(id)
);

