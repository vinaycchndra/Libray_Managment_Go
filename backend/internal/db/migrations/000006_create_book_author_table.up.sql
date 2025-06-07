CREATE TABLE book_borrow (
    id SERIAL primary key,  
    book_id INTEGER NOT NULL, 
    list_id INTEGER NOT NULL,
    returned BOOLEAN DEFAULT false, 
    fine_paid NUMERIC(6, 2) DEFAULT 0, 
    extended BOOLEAN DEFAULT false, 
    FOREIGN KEY (book_id) REFERENCES book(id), 
    FOREIGN KEY (list_id) REFERENCES book_borrow_list(id)
);