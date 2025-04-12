CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY, 
    name VARCHAR(64) NOT NULL, 
    email VARCHAR(32) NOT NULL UNIQUE, 
    password   VARCHAR(256) NOT NULL, 
    phone_number VARCHAR(16) NOT NULL, 
    created_at TIMESTAMP NOT NULL, 
    updated_at TIMESTAMP NOT NULL, 
    is_active BOOLEAN DEFAULT false, 
    is_admin BOOLEAN DEFAULT false
);
