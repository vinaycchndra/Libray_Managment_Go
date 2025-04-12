CREATE TABLE book (
    id SERIAL PRIMARY KEY, 
    title VARCHAR(256) NOT NULL, 
    category VARCHAR(64) NOT NULL,
    publisher VARCHAR(64) NOT NULL, 
    book_count INTEGER DEFAULT 0, 
    price NUMERIC(10, 2), 
    fine_per_day NUMERIC(6, 2),
    author_id INTEGER NOT NULL, 
    created_at TIMESTAMP NOT NULL, 
    updated_at TIMESTAMP NOT NULL, 
    archive BOOLEAN DEFAULT false, 
    FOREIGN KEY (author_id) REFERENCES author(id)
); 