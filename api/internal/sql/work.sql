CREATE TABLE work
(
    id SERIAL PRIMARY KEY,
    author_id INTEGER NOT NULL,
    description VARCHAR (5000) NOT NULL,
    initial_pub_date DATE NOT NULL,
    original_language VARCHAR (100) NOT NULL,
    title VARCHAR (200) NOT NULL,
    UNIQUE (author_id, title),
    FOREIGN KEY (author_id) REFERENCES author (id)
);
