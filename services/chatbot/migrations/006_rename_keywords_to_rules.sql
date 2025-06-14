-- Rename main table and all related tables
ALTER TABLE keywords RENAME TO rules;
ALTER TABLE keyword_responses RENAME TO rule_responses;
ALTER TABLE keyword_pages RENAME TO rule_pages;

-- Rename columns
ALTER TABLE rule_responses RENAME COLUMN keyword_id TO rule_id;
ALTER TABLE rule_pages RENAME COLUMN keyword_id TO rule_id;
ALTER TABLE rate_limits RENAME COLUMN keyword_id TO rule_id;
ALTER TABLE incoming_messages RENAME COLUMN matched_keyword_id TO matched_rule_id;
ALTER TABLE rule_analytics RENAME COLUMN keyword_id TO rule_id;

-- Update sequences
ALTER SEQUENCE keywords_id_seq RENAME TO rules_id_seq;
