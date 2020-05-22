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

INSERT INTO author
    (first_name, last_name, gender, date_of_birth, place_of_birth)
VALUES
    (
        'Sally',
        'Rooney',
        'F',
        '1991-02-20',
        (
            SELECT id
            FROM location
            WHERE city = 'Castlebar' AND country = 'Ireland'
        )
    ),
    (
        'Celeste',
        'Ng',
        'F',
        '1980-07-30',
        (
            SELECT id
            FROM location
            WHERE city = 'Pittsburgh' AND country = 'USA'
        )
    ),
    (
        'Delia',
        'Owens',
        'F',
        '1949-01-01',
        (
            SELECT id
            FROM location
            WHERE city = 'Thomasville' AND country = 'USA'
        )
    ),
    (
        'F. Scott',
        'Fitzgerald',
        'M',
        '1896-09-24',
        (
            SELECT id
            FROM location
            WHERE city = 'Saint Paul' AND country = 'USA'
        )
    ),
    (
        'Jeanine',
        'Cummins',
        'F',
        NULL,
        (
            SELECT id
            FROM location
            WHERE city = 'Rota' AND country = 'Spain'
        )
    ),
    (
        'James',
        'Clear',
        'M',
        NULL,
        (
            SELECT id
            FROM location
            WHERE city = 'Hamilton' AND country = 'USA'
        )
    );
