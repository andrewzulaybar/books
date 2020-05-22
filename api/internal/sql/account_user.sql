CREATE TABLE account_user
(
    id SERIAL PRIMARY KEY,
    email VARCHAR (100) UNIQUE NOT NULL,
    first_name VARCHAR (100) NOT NULL,
    last_name VARCHAR (100) NOT NULL,
    gender VARCHAR (1) NOT NULL,
    password VARCHAR (200) NOT NULL,
    username VARCHAR (100) UNIQUE NOT NULL,
    age INTEGER,
    location_id INTEGER,
    time_created TIMESTAMP NOT NULL,
    FOREIGN KEY (location_id) REFERENCES location (id)
);

INSERT INTO account_user
    (email, first_name, last_name, gender, password, username, age, location_id, time_created)
VALUES
    (
        'andrewzulaybar@gmail.com',
        'Andrew',
        'Zulaybar',
        'M',
        'password',
        'admin',
        '21',
        (
            SELECT id
            FROM location
            WHERE city = 'Vancouver' AND country = 'Canada'
        ),
        CURRENT_TIMESTAMP
    );

