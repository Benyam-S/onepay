CREATE TABLE users (
    user_id VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    phone_number VARCHAR(255) NOT NULL
);