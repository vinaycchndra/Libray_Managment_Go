CREATE TABLE author (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL, 
    about TEXT NOT NULL, 
    created_at TIMESTAMP NOT NULL, 
    updated_at TIMESTAMP NOT NULL 
);