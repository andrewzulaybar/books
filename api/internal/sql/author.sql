CREATE TABLE author
(
    id SERIAL PRIMARY KEY,
    first_name VARCHAR (100) NOT NULL,
    last_name VARCHAR (100) NOT NULL,
    gender VARCHAR (1) NOT NULL,
    date_of_birth DATE NOT NULL,
    place_of_birth INTEGER NOT NULL,
    UNIQUE (first_name, last_name, date_of_birth),
    FOREIGN KEY (place_of_birth) REFERENCES location (id)
);
