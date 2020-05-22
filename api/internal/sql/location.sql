CREATE TABLE location
(
    id SERIAL PRIMARY KEY,
    city VARCHAR (100) NOT NULL,
    country VARCHAR (100) NOT NULL,
    region VARCHAR (100) NOT NULL,
    UNIQUE (city, country)
);

INSERT INTO location
    (city, region, country)
VALUES
    ('Vancouver', 'British Columbia', 'Canada'),
    ('Toronto', 'Ontario', 'Canada'),
    ('Montreal', 'Quebec', 'Canada'),
    ('Castlebar', 'County Mayo', 'Ireland'),
    ('Rota', 'Andalusia', 'Spain'),
    ('ThomasVille', 'Georgia', 'USA'),
    ('Saint Paul', 'Minnesota', 'USA'),
    ('New York City', 'New York', 'USA'),
    ('Hamilton', 'Ohio', 'USA'),
    ('Pittsburgh', 'Pensylvannia', 'USA');

