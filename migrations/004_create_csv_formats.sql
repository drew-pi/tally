CREATE TABLE IF NOT EXISTS csv_formats (
    id          SERIAL PRIMARY KEY,
    bank_id     INTEGER REFERENCES banks(id),
    csv_column  TEXT NOT NULL,
    column_type TEXT NOT NULL
);