CREATE TYPE card_type AS ENUM ('WF', 'Fidelity');

CREATE TABLE IF NOT EXISTS transactions (
    id          SERIAL PRIMARY KEY,
    date        DATE,
    description TEXT,
    amount      NUMERIC(10, 2),
    card        card_type
);