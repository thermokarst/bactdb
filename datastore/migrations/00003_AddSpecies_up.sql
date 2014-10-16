-- bactdb
-- Matthew Dillon

CREATE TABLE species (
    id BIGSERIAL NOT NULL,
    genusid BIGINT NOT NULL,
    speciesname CHARACTER VARYING(100),

    createdat TIMESTAMP WITH TIME ZONE,
    updatedat TIMESTAMP WITH TIME ZONE,
    deletedat TIMESTAMP WITH TIME ZONE,

    CONSTRAINT species_pkey PRIMARY KEY (id),
    FOREIGN KEY (genusid) REFERENCES genera(id)
);
