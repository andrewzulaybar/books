CREATE TABLE location
(
    id SERIAL PRIMARY KEY,
    city VARCHAR (100) NOT NULL,
    country VARCHAR (100) NOT NULL,
    region VARCHAR (100) NOT NULL,
    UNIQUE (city, country)
);

