CREATE TABLE book_borrow (
    id SERIAL primary key, 
    user_id INTEGER NOT NULL, 
    book_id INTEGER NOT NULL, 
    due_date TIMESTAMP NOT NULL, 
    returned BOOLEAN DEFAULT false, 
    fine_paid NUMERIC(6, 2) DEFAULT 0, 
    extended BOOLEAN DEFAULT false, 
    created_at TIMESTAMP NOT NULL, 
    updated_at TIMESTAMP NOT NULL, 
    FOREIGN KEY (book_id) REFERENCES book(id), 
    FOREIGN KEY (user_id) REFERENCES users(id)
);