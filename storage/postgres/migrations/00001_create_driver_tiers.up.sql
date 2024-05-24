CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE tiers (
    id uuid DEFAULT uuid_generate_v4(),
    frequency float,
    last_name text,
    PRIMARY KEY(id)
);

-- demo purposes
INSERT INTO fighters (first_name, last_name) VALUES ('Dave', 'Grohl');