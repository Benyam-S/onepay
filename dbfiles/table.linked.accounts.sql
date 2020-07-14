CREATE TABLE linked_accounts(
    id VARCHAR PRIMARY KEY,
	user_id VARCHAR,
    account_provider VARCHAR,
	account_id VARCHAR,
	access_token VARCHAR
);