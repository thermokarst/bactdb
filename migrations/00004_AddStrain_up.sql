-- bactdb
-- Matthew R Dillon

CREATE TABLE strains (
    id BIGSERIAL NOT NULL,
    species_id BIGINT NOT NULL,
    strain_name TEXT NOT NULL,
    type_strain BOOLEAN NOT NULL,
    accession_numbers TEXT NULL,
    genbank TEXT NULL,
    whole_genome_sequence TEXT NULL,
    isolated_from TEXT NULL,
    notes TEXT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    created_by BIGINT NOT NULL,
    updated_by BIGINT NOT NULL,

    CONSTRAINT strain_pkey PRIMARY KEY (id),
    FOREIGN KEY (species_id) REFERENCES species(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
);

CREATE INDEX species_id_idx ON strains (species_id);

