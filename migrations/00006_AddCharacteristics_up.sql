-- bactdb
-- Matthew R Dillon

CREATE TABLE characteristics (
    id BIGSERIAL NOT NULL,
    characteristic_name TEXT NOT NULL,
    characteristic_type_id BIGINT NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    CONSTRAINT characteristics_pkey PRIMARY KEY (id),
    FOREIGN KEY (characteristic_type_id) REFERENCES characteristic_types(id)
);

CREATE INDEX characteristic_type_id_idx ON characteristics (characteristic_type_id);

