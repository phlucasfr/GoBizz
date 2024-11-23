DROP INDEX IF EXISTS idx_company_email_active;
DROP INDEX IF EXISTS idx_company_phone_active;

DROP INDEX IF EXISTS idx_company_name;
DROP INDEX IF EXISTS idx_company_email;
DROP INDEX IF EXISTS idx_company_phone;
DROP INDEX IF EXISTS idx_company_updated_at;

DROP TABLE IF EXISTS "company";

DROP EXTENSION IF EXISTS "uuid-ossp";
