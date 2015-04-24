-- bactdb
-- Matthew Dillon

CREATE TABLE species (
    id BIGSERIAL NOT NULL,
    genus_id BIGINT NOT NULL,
    species_name CHARACTER VARYING(100) NOT NULL,
    type_species BOOLEAN NULL,
    subspecies_species_id BIGINT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    CONSTRAINT species_pkey PRIMARY KEY (id),
    FOREIGN KEY (genus_id) REFERENCES genera(id),
    FOREIGN KEY (subspecies_species_id) REFERENCES species(id)
);

CREATE INDEX genus_id_idx ON species (genus_id);

