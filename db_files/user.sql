CREATE TABLE user (
    user_id VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255),
    phone_number VARCHAR(255)
);

CREATE TABLE password(
    user_id VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL,
    password VARCHAR(255),
    salt VARCHAR(255)
)