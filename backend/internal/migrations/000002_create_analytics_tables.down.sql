-- Drop analytics tables
DROP INDEX IF EXISTS idx_webhook_logs_status;
DROP INDEX IF EXISTS idx_webhook_logs_created_at;
DROP INDEX IF EXISTS idx_webhook_logs_type;
DROP TABLE IF EXISTS webhook_logs;

DROP INDEX IF EXISTS idx_analytics_created_at;
DROP INDEX IF EXISTS idx_analytics_keyword;
DROP INDEX IF EXISTS idx_analytics_recipient;
DROP INDEX IF EXISTS idx_analytics_sender;
DROP TABLE IF EXISTS message_analytics;
