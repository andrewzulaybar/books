DROP TABLE IF EXISTS publication;

CREATE TABLE publication
(
    id SERIAL PRIMARY KEY,
    author VARCHAR (100) NOT NULL,
    edition_pub_date DATE NOT NULL,
    format VARCHAR (30) NOT NULL,
    image_url VARCHAR (300) UNIQUE NOT NULL,
    isbn VARCHAR (10) UNIQUE NOT NULL,
    isbn13 VARCHAR (13) UNIQUE NOT NULL,
    language VARCHAR (100) NOT NULL,
    num_pages INTEGER NOT NULL,
    publisher VARCHAR (100) NOT NULL,
    title VARCHAR (200) NOT NULL,
    UNIQUE (author, title)
);
