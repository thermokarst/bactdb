-- bactdb
-- Matthew R Dillon

CREATE TABLE characteristics (
    id BIGSERIAL NOT NULL,
    characteristic_name TEXT NOT NULL,
    characteristic_type_id BIGINT NOT NULL,
    sort_order BIGINT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    created_by BIGINT NOT NULL,
    updated_by BIGINT NOT NULL,
    deleted_by BIGINT NULL,

    CONSTRAINT characteristics_pkey PRIMARY KEY (id),
    FOREIGN KEY (characteristic_type_id) REFERENCES characteristic_types(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id),
    FOREIGN KEY (deleted_by) REFERENCES users(id)
);

CREATE INDEX characteristic_type_id_idx ON characteristics (characteristic_type_id);

