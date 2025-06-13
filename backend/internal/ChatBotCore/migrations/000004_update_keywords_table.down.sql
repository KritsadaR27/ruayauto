-- Revert keywords table changes
DROP INDEX IF EXISTS idx_keywords_keyword_text;
DROP INDEX IF EXISTS idx_keywords_active_priority;
DROP INDEX IF EXISTS idx_keywords_active;
DROP INDEX IF EXISTS idx_keywords_priority;
DROP INDEX IF EXISTS idx_keywords_match_type;

-- Remove added columns
ALTER TABLE keywords 
DROP COLUMN IF EXISTS active,
DROP COLUMN IF EXISTS priority,
DROP COLUMN IF EXISTS match_type;
