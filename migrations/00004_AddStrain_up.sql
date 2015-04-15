-- bactdb
-- Matthew R Dillon

CREATE TABLE strains (
    id BIGSERIAL NOT NULL,
    species_id BIGINT NOT NULL,
    strain_name CHARACTER VARYING(100) NOT NULL,
    type_strain BOOLEAN NOT NULL,
    etymology CHARACTER VARYING(500) NULL,
    accession_banks CHARACTER VARYING(100) NULL,
    genbank_embl_ddb CHARACTER VARYING(100) NULL,
    isolated_from CHARACTER VARYING(100) NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    CONSTRAINT strain_pkey PRIMARY KEY (id),
    FOREIGN KEY (species_id) REFERENCES species(id)
);

CREATE INDEX species_id_idx ON strains (species_id);

