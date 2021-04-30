CREATE TABLE IF NOT EXISTS book (
    id serial PRIMARY KEY,
    isbn VARCHAR ( 255 ) NOT NULL,
    name VARCHAR ( 255 ) NOT NULL,
    image VARCHAR ( 255 ) NOT NULL,
    genre VARCHAR ( 255 ) NOT NULL,
    year_published int NOT NULL,
);