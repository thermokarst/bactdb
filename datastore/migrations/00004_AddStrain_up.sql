-- bactdb
-- Matthew R Dillon

CREATE TABLE strains (
    id BIGSERIAL NOT NULL,
    species_id BIGINT,
    strain_name CHARACTER VARYING(100),
    strain_type CHARACTER VARYING(100),
    etymology CHARACTER VARYING(500),
    accession_banks CHARACTER VARYING(100),
    genbank_embl_ddb CHARACTER VARYING(100),

    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT strain_pkey PRIMARY KEY (id),
    FOREIGN KEY (species_id) REFERENCES species(id)
);

