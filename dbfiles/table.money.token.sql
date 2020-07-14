CREATE TABLE money_tokens(
    code VARCHAR PRIMARY KEY,
    sender_id VARCHAR,
    sent_at DATETIME,
    amount FLOAT,
    expiration_date DATETIME,
    method VARCHAR
);