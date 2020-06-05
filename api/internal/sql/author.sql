CREATE TABLE author
(
    id SERIAL PRIMARY KEY,
    first_name VARCHAR (100) NOT NULL,
    last_name VARCHAR (100) NOT NULL,
    gender VARCHAR (1) NOT NULL,
    date_of_birth DATE,
    place_of_birth INTEGER,
    FOREIGN KEY (place_of_birth) REFERENCES location (id)
);
