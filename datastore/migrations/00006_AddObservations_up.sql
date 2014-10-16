-- bactdb
-- Matthew R Dillon

CREATE TABLE observations (
    id BIGSERIAL NOT NULL,
    observation_name CHARACTER VARYING(100),
    observation_type_id BIGINT,

    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT observations_pkey PRIMARY KEY (id),
    FOREIGN KEY (observation_type_id) REFERENCES observation_types(id)
);

