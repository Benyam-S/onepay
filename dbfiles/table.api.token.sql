CREATE TABLE api_tokens(
    api_token VARCHAR UNIQUE,
    api_key VARCHAR, -- can be used to identify the app
    user_id VARCHAR, -- for whom it is created
    expires_at INT,
    daily_expiration INT,
);