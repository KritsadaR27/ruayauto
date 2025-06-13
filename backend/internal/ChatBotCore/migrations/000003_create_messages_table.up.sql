-- Create messages table for storing all conversation messages
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    conversation_id INTEGER NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    message_id VARCHAR(255) UNIQUE, -- Facebook message ID (nullable for bot messages)
    sender_type VARCHAR(10) NOT NULL CHECK (sender_type IN ('user', 'bot')),
    content TEXT,
    message_type VARCHAR(20) DEFAULT 'text' CHECK (message_type IN ('text', 'image', 'sticker', 'video', 'audio', 'file')),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP WITH TIME ZONE,
    response_time_ms INTEGER -- Time taken to generate response
);

-- Create indexes for better query performance
CREATE INDEX idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX idx_messages_sender_type ON messages(sender_type);
CREATE INDEX idx_messages_message_type ON messages(message_type);
CREATE INDEX idx_messages_created_at ON messages(created_at DESC);
CREATE INDEX idx_messages_message_id ON messages(message_id) WHERE message_id IS NOT NULL;

-- Create index on metadata for JSON queries
CREATE INDEX idx_messages_metadata_gin ON messages USING gin(metadata);

-- Add foreign key constraint with proper naming
ALTER TABLE messages 
ADD CONSTRAINT fk_messages_conversation_id 
FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE;
