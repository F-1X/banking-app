CREATE SCHEMA IF NOT EXISTS payment_system;

CREATE TABLE IF NOT EXISTS payment_system.wallets (
    id SERIAL PRIMARY KEY,
    address VARCHAR(256) NOT NULL,
    balance REAL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS payment_system.transactions (
    id SERIAL PRIMARY KEY,
    from_address VARCHAR(256) NOT NULL,
    to_address VARCHAR(256) NOT NULL,
    amount REAL
);

INSERT INTO payment_system.wallets (address, balance)
SELECT * FROM (
    SELECT '1' AS address, 100 AS balance UNION ALL
    SELECT '2', 100 UNION ALL
    SELECT '3', 100 UNION ALL
    SELECT '4', 100 UNION ALL
    SELECT '5', 100 UNION ALL
    SELECT '6', 100 UNION ALL
    SELECT '7', 100 UNION ALL
    SELECT '8', 100 UNION ALL
    SELECT '9', 100 UNION ALL
    SELECT '10', 100
) AS tmp
WHERE NOT EXISTS (SELECT 1 FROM payment_system.wallets);
