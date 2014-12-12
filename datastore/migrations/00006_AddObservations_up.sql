-- bactdb
-- Matthew R Dillon

CREATE TABLE observations (
    id BIGSERIAL NOT NULL,
    observation_name CHARACTER VARYING(100) NOT NULL,
    characteristic_type_id BIGINT NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    CONSTRAINT observations_pkey PRIMARY KEY (id),
    FOREIGN KEY (characteristic_type_id) REFERENCES characteristic_types(id)
);

CREATE INDEX characteristic_type_id_idx ON observations (characteristic_type_id);

