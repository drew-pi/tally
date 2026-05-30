CREATE TYPE payment_method_type AS ENUM ('credit', 'debit', 'cheque', 'transfer');

CREATE TABLE IF NOT EXISTS payment_methods (
    id      SERIAL PRIMARY KEY,
    type    payment_method_type NOT NULL,
    bank_id INTEGER REFERENCES banks(id)
);