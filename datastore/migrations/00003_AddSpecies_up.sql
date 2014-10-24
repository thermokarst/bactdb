-- bactdb
-- Matthew Dillon

CREATE TABLE species (
    id BIGSERIAL NOT NULL,
    genus_id BIGINT NOT NULL,
    species_name CHARACTER VARYING(100),

    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT species_pkey PRIMARY KEY (id),
    FOREIGN KEY (genus_id) REFERENCES genera(id)
);
