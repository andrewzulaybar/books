DROP TABLE IF EXISTS work, publication;

CREATE TABLE work
(
    id SERIAL PRIMARY KEY,
    author VARCHAR (100) NOT NULL,
    description VARCHAR (5000) NOT NULL,
    initial_pub_date DATE NOT NULL,
    original_language VARCHAR (100) NOT NULL,
    title VARCHAR (200) NOT NULL,
    UNIQUE (author, title)
);

CREATE TABLE publication
(
    id SERIAL PRIMARY KEY,
    edition_pub_date DATE NOT NULL,
    format VARCHAR (30) NOT NULL,
    image_url VARCHAR (300) UNIQUE NOT NULL,
    isbn VARCHAR (10) UNIQUE NOT NULL,
    isbn13 VARCHAR (13) UNIQUE NOT NULL,
    language VARCHAR (100) NOT NULL,
    num_pages INTEGER NOT NULL,
    publisher VARCHAR (100) NOT NULL,
    work_id INTEGER NOT NULL,
    FOREIGN KEY (work_id) REFERENCES work (id)
);
