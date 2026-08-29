DROP TABLE IF EXISTS business_claims;
ALTER TABLE businesses DROP CONSTRAINT IF EXISTS businesses_owner_id_fkey;
