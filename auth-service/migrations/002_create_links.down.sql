-- Drop trigger and function
DROP TRIGGER IF EXISTS update_links_updated_at ON links;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_links_short_url;
DROP INDEX IF EXISTS idx_links_custom_slug;
DROP INDEX IF EXISTS idx_links_customer_id;

-- Drop table
DROP TABLE IF EXISTS links; 