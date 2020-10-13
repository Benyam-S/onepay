CREATE TABLE deleted_linked_accounts(
    id VARCHAR PRIMARY KEY,
    user_id VARCHAR,
    account_provider_id VARCHAR,
    account_id VARCHAR,
    access_token VARCHAR
);