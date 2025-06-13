-- Drop conversations table and related objects
DROP TRIGGER IF EXISTS trigger_conversations_updated_at ON conversations;
DROP FUNCTION IF EXISTS update_conversations_updated_at();
DROP INDEX IF EXISTS idx_conversations_unique_user_post;
DROP INDEX IF EXISTS idx_conversations_last_message_at;
DROP INDEX IF EXISTS idx_conversations_status;
DROP INDEX IF EXISTS idx_conversations_source_type;
DROP INDEX IF EXISTS idx_conversations_facebook_user_id;
DROP INDEX IF EXISTS idx_conversations_page_id;
DROP TABLE IF EXISTS conversations;
