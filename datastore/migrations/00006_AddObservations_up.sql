-- bactdb
-- Matthew R Dillon

CREATE TABLE observations (
    id BIGSERIAL NOT NULL,
    observation_name CHARACTER VARYING(100) NOT NULL,
    observation_type_id BIGINT NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    CONSTRAINT observations_pkey PRIMARY KEY (id),
    FOREIGN KEY (observation_type_id) REFERENCES observation_types(id)
);

CREATE INDEX observation_type_id_idx ON observations (observation_type_id);

