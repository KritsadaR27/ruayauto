-- Drop keywords table and related objects
DROP TRIGGER IF EXISTS update_keywords_updated_at ON keywords;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_keywords_active;
DROP INDEX IF EXISTS idx_keywords_keyword;
DROP TABLE IF EXISTS keywords;
