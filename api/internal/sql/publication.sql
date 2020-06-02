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
    FOREIGN KEY (work_id) REFERENCES work (id) ON DELETE CASCADE
);

INSERT INTO publication
    (edition_pub_date, format, image_url, isbn, isbn13, language, num_pages, publisher, work_id)
VALUES
    (
        '2019-04-16',
        'Hardcover',
        'https://images-na.ssl-images-amazon.com/images/I/81X4R7QhFkL.jpg',
        '1984822179',
        '9781984822178',
        'English',
        288,
        'Hogarth',
        (
            SELECT id
            FROM work
            WHERE author_id =
                (
                    SELECT id
                    FROM author
                    WHERE first_name = 'Sally' AND last_name = 'Rooney'
                )
                AND title = 'Normal People'
        )
    ),
    (
        '2017-09-12',
        'Hardcover',
        'https://images-na.ssl-images-amazon.com/images/I/91twTG-CQ8L.jpg',
        '0735224293',
        '9780735224292',
        'English',
        352,
        'Penguin Press',
        (
            SELECT id
            FROM work
            WHERE author_id =
                (
                    SELECT id
                    FROM author
                    WHERE first_name = 'Celeste' AND last_name = 'Ng'
                )
                AND title = 'Little Fires Everywhere'
        )
    ),
    (
        '2018-08-14',
        'Hardcover',
        'https://images-na.ssl-images-amazon.com/images/I/51j5p18mJNL.jpg',
        '0735219095',
        '9780735219090',
        'English',
        384,
        'G.P. Putnam''s Sons',
        (
            SELECT id
            FROM work
            WHERE author_id =
                (
                    SELECT id
                    FROM author
                    WHERE first_name = 'Delia' AND last_name = 'Owens'
                )
                AND title = 'Where the Crawdads Sing'
        )
    ),
    (
        '2004-09-30',
        'Paperback',
        'https://images-na.ssl-images-amazon.com/images/I/81af+MCATTL.jpg',
        '0743273567',
        '9780743273565',
        'English',
        180,
        'Scribner',
        (
            SELECT id
            FROM work
            WHERE author_id =
                (
                    SELECT id
                    FROM author
                    WHERE first_name = 'F. Scott' AND last_name = 'Fitzgerald'
                )
                AND title = 'The Great Gatsby'
        )
    ),
    (
        '2020-01-21',
        'Hardcover',
        'https://images-na.ssl-images-amazon.com/images/I/81iVsj91eQL.jpg',
        '1250209765',
        '9781250209764',
        'English',
        400,
        'Flatiron Books',
        (
            SELECT id
            FROM work
            WHERE author_id =
                (
                    SELECT id
                    FROM author
                    WHERE first_name = 'Jeanine' AND last_name = 'Cummins'
                )
                AND title = 'American Dirt'
        )
    ),
    (
        '2018-09-16',
        'Hardcover',
        'https://images-na.ssl-images-amazon.com/images/I/91Xq+S+F2jL.jpg',
        '0735211299',
        '9780735211292',
        'English',
        320,
        'Avery',
        (
            SELECT id
            FROM work
            WHERE author_id =
                (
                    SELECT id
                    FROM author
                    WHERE first_name = 'James' AND last_name = 'Clear'
                )
                AND title = 'Atomic Habits'
        )
    );
