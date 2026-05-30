CREATE TABLE IF NOT EXISTS transactions (
    id                SERIAL PRIMARY KEY,
    date              DATE          NOT NULL,
    vendor            TEXT          NOT NULL,
    description       TEXT,
    category          TEXT,
    amount            NUMERIC(10,2) NOT NULL,
    payment_method_id INTEGER REFERENCES payment_methods(id)
);