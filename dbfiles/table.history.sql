CREATE TABLE user_history (
    id INT PRIMARY KEY,
    sender_id VARCHAR,
    receiver_id VARCHAR,
    sent_at VARCHAR,
    received_at VARCHAR,
    method VARCHAR,
    code VARCHAR,
    amount FLOAT,
    sender_seen BOOLEAN,
    receiver_seen BOOLEAN
);