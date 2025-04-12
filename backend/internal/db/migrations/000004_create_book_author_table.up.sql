CREATE TABLE book_author (
    id SERIAL PRIMARY KEY, 
    book_id INTEGER NOT NULL, 
    author_id INTEGER NOT NULL, 
    created_at TIMESTAMP NOT NULL, 
    updated_at TIMESTAMP NOT NULL, 
    UNIQUE (author_id, book_id), 
    FOREIGN KEY (author_id) REFERENCES author(id), 
    FOREIGN KEY (book_id) REFERENCES book(id)
); 