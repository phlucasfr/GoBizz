DROP INDEX IF EXISTS idx_customer_email_active;
DROP INDEX IF EXISTS idx_customer_phone_active;

DROP INDEX IF EXISTS idx_customer_name;
DROP INDEX IF EXISTS idx_customer_email;
DROP INDEX IF EXISTS idx_customer_phone;
DROP INDEX IF EXISTS idx_customer_updated_at;

DROP TABLE IF EXISTS "customer";

DROP EXTENSION IF EXISTS "uuid-ossp";