CREATE TABLE users (
    user_id VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    phone_number VARCHAR(255) NOT NULL
);

CREATE TABLE user_passwords(
    user_id VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    salt VARCHAR(255) NOT NULL
);