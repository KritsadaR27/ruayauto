-- Drop messages table and related objects
ALTER TABLE messages DROP CONSTRAINT IF EXISTS fk_messages_conversation_id;
DROP INDEX IF EXISTS idx_messages_metadata_gin;
DROP INDEX IF EXISTS idx_messages_message_id;
DROP INDEX IF EXISTS idx_messages_created_at;
DROP INDEX IF EXISTS idx_messages_message_type;
DROP INDEX IF EXISTS idx_messages_sender_type;
DROP INDEX IF EXISTS idx_messages_conversation_id;
DROP TABLE IF EXISTS messages;
