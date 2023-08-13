-- drop the constraints 
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";

-- drop index
DROP INDEX IF EXISTS "accounts_owner_currency_idx";

DROP TABLE IF EXISTS "users";