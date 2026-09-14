-- Uniform seed data.
-- Pattern: DROP + CREATE + INSERT so every dev boot matches this file exactly.
-- New items: append as the next ID, don't renumber.

DROP TABLE IF EXISTS uniform;

CREATE TABLE uniform (
    id         INTEGER PRIMARY KEY NOT NULL,
    name       TEXT    UNIQUE NOT NULL,
    weight_lbs REAL    NOT NULL
);

INSERT INTO uniform (id, name, weight_lbs) VALUES
    (1, 'Tropical', 4.5),
    (2, 'Normal',   5),
    (2, 'Winter',   7);