-- Migration: Split keywords from rules table to rule_keywords table (1 rule : N keyword)

-- 1. Create rule_keywords table
CREATE TABLE IF NOT EXISTS rule_keywords (
    id SERIAL PRIMARY KEY,
    rule_id INTEGER NOT NULL REFERENCES rules(id) ON DELETE CASCADE,
    keyword VARCHAR(255) NOT NULL
);

-- 2. Migrate data from rules.keyword to rule_keywords (split by comma if needed)
INSERT INTO rule_keywords (rule_id, keyword)
SELECT id, trim(keyword) FROM rules WHERE keyword IS NOT NULL AND length(trim(keyword)) > 0;

-- 3. Drop dependent view before dropping column
DROP VIEW IF EXISTS keywords_with_responses;

-- 4. Drop keyword column from rules table
ALTER TABLE rules DROP COLUMN IF EXISTS keyword;
