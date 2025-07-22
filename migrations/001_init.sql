-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    UNIQUE(username)
);

INSERT INTO users (username, password) VALUES ('admin', '12345');

CREATE TABLE IF NOT EXISTS clients (
    user_id INT,
    name VARCHAR(100) NOT NULL,
    passport VARCHAR(100) NOT NULL,
    usd FLOAT,
    eur FLOAT,
    currency VARCHAR(3) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT clock_timestamp(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT clock_timestamp(),
    UNIQUE(passport),
    CONSTRAINT fk_user_clients
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON UPDATE CASCADE
    ON DELETE SET NULL
    );

CREATE TABLE IF NOT EXISTS client_histories (
    user_id INT,
    passport VARCHAR(100) NOT NULL,
    usd FLOAT,
    eur FLOAT,
    currency VARCHAR(3) NOT NULL,
    is_plus BOOLEAN NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT clock_timestamp(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT clock_timestamp(),
    CONSTRAINT fk_user_histories
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON UPDATE CASCADE
    ON DELETE SET NULL
    );

-- +goose Down
drop table users;
drop table clients;
drop table client_histories;
