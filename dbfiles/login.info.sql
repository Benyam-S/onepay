Create table server_sessions(
    id INT PRIMARY KEY UNIQUE NOT NULL AUTO_INCREMENT,
    user_id VARCHAR NOT NULL,
    session_id VARCHAR UNIQUE NOT NULL,
    ip_address VARCHAR,
    device_info VARCHAR
);