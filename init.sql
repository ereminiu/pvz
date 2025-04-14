CREATE TABLE users (
    id INT PRIMARY KEY
);

INSERT INTO users (id) 
    VALUES (19), (26), (12);

CREATE TABLE orders (
    id INT PRIMARY KEY,
    user_id INT REFERENCES users(id),
    weight INT NOT NULL,
    price INT NOT NULL,
    packing VARCHAR(10) NOT NULL,
    extra BOOLEAN NOT NULL,
    status varchar(20) NOT NULL,
    expire_at TIMESTAMP NOT NULL,
    created_moment SERIAL,
    created_at TIMESTAMP NOT NULL DEFAULT(NOW()),
    updated_at TIMESTAMP NOT NULL DEFAULT(NOW())
);

CREATE TABLE admins (
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE logs (
    id SERIAL PRIMARY KEY,
    user_id int,
    action VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT(NOW()),
    error TEXT NOT NULL
);