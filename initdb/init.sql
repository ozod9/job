create type transaction_type as enum ('income', 'outcome');

DROP TABLE IF EXISTS balances CASCADE;
DROP TABLE IF EXISTS transactions;

CREATE TABLE IF NOT EXISTS balances
(
   	ID SERIAL PRIMARY KEY,
	balance DECIMAL CHECK (balance >= 0)
);


CREATE TABLE IF NOT EXISTS transactions
(
    ID SERIAL PRIMARY KEY,
	balance_id INTEGER,
    from_id INTEGER,
    amount DECIMAL CHECK (amount >= 0),
	reason CHARACTER VARYING(50) NOT NULL,
	type transaction_type NOT NULL,
	date timestamptz NOT NULL, 
    FOREIGN KEY(balance_id) REFERENCES balances(id) ON DELETE CASCADE
);


SELECT * FROM balances;