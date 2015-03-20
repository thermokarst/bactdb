-- bactdb
-- Matthew R Dillon

CREATE TABLE characteristic_types (
    id BIGSERIAL NOT NULL,
    characteristic_type_name CHARACTER VARYING(100) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    CONSTRAINT characteristic_types_pkey PRIMARY KEY (id)
);

