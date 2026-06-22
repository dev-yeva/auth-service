CREATE TABLE IF NOT EXISTS users
(
    id                 INTEGER PRIMARY KEY,
    email              TEXT NOT NULL UNIQUE,
    is_admin           BOOLEAN NOT NULL DEFAULT 0,
    password_hash BLOB NOT NULL
);

INSERT INTO users (email, is_admin, password_hash) 
VALUES ('boss@gmail.com', true, '$2a$04$zr08lx1iof//9imZMh47h.9O3c0frMP0D4nzH6NUEFvUCYgQTPyAC');

CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE TABLE IF NOT EXISTS apps
(
    id         INTEGER PRIMARY KEY,
    title      TEXT NOT NULL,
    secret_key TEXT NOT NULL UNIQUE
);

INSERT INTO apps (title, secret_key) 
VALUES 
    ('Gemini', '5NqgcSG0YFyz7YHAVPMC'),
    ('Gmail', '9blmaMXKfu19KSqTvd8f'),
    ('Chrome', 'sZmV3ABYqbSJqY0cs73F');