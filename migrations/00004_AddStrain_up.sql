-- bactdb
-- Matthew R Dillon

CREATE TABLE strains (
    id BIGSERIAL NOT NULL,
    species_id BIGINT NOT NULL,
    strain_name TEXT NOT NULL,
    type_strain BOOLEAN NOT NULL,
    accession_banks TEXT NULL,
    genbank_embl_ddb TEXT NULL,
    isolated_from TEXT NULL,
    notes TEXT NULL,
    author_id BIGINT NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    CONSTRAINT strain_pkey PRIMARY KEY (id),
    FOREIGN KEY (species_id) REFERENCES species(id),
    FOREIGN KEY (author_id) REFERENCES users(id)
);

CREATE INDEX species_id_idx ON strains (species_id);

