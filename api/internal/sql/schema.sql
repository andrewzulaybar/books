DROP TABLE IF EXISTS publication;

CREATE TABLE publication
(
    id SERIAL PRIMARY KEY,
    author VARCHAR (100) NOT NULL,
    image_url VARCHAR (300) NOT NULL,
    title VARCHAR (200) NOT NULL
);
