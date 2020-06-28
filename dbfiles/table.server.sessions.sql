Create table server_sessions(
    user_id VARCHAR NOT NULL,
    session_id VARCHAR UNIQUE NOT NULL,
    ip_address VARCHAR,
    device_info VARCHAR
);