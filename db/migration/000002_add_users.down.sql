ALTER TABLE if exists "accounts" drop constraint if exists "owner_currency_key";
ALTER TABLE if exists "accounts" drop constraint if exists "accounts_owner_fkey";
DROP TABLE IF EXISTS users;
