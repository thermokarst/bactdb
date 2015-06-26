-- bactdb
-- Matthew R Dillon

CREATE TABLE verification (
    user_id BIGINT NOT NULL,
    nonce CHARACTER(60) NOT NULL UNIQUE,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,

    CONSTRAINT verification_pkey PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

