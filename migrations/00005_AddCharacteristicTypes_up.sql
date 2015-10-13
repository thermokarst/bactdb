-- bactdb
-- Matthew R Dillon

CREATE TABLE characteristic_types (
    id BIGSERIAL NOT NULL,
    characteristic_type_name TEXT NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    created_by BIGINT NOT NULL,
    updated_by BIGINT NOT NULL,

    CONSTRAINT characteristic_types_pkey PRIMARY KEY (id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
);

