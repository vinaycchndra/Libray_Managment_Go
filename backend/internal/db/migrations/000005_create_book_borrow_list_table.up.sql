CREATE TABLE book_borrow_list (
    id SERIAL primary key,    
    due_date TIMESTAMP NOT NULL,
    user_id INTEGER NOT NULL, 
    closed BOOLEAN DEFAULT false, 
    fine_paid NUMERIC(6, 2) DEFAULT 0,  
    created_at TIMESTAMP NOT NULL, 
    updated_at TIMESTAMP NOT NULL,  
    FOREIGN KEY (user_id) REFERENCES users(id)
);