CREATE TABLE IF NOT EXISTS banks (
    id   SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

INSERT INTO banks (name) VALUES ('WF'), ('Fidelity');