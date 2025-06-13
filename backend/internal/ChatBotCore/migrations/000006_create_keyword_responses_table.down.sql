-- Drop keyword_responses table and related objects
ALTER TABLE keyword_responses DROP CONSTRAINT IF EXISTS fk_keyword_responses_response_template_id;
ALTER TABLE keyword_responses DROP CONSTRAINT IF EXISTS fk_keyword_responses_keyword_id;
DROP INDEX IF EXISTS idx_keyword_responses_active_order;
DROP INDEX IF EXISTS idx_keyword_responses_active;
DROP INDEX IF EXISTS idx_keyword_responses_order_index;
DROP INDEX IF EXISTS idx_keyword_responses_response_template_id;
DROP INDEX IF EXISTS idx_keyword_responses_keyword_id;
DROP TABLE IF EXISTS keyword_responses;
