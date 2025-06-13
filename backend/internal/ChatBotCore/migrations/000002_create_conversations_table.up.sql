-- Create conversations table for tracking Facebook conversations
CREATE TABLE IF NOT EXISTS conversations (
    id SERIAL PRIMARY KEY,
    page_id VARCHAR(255) NOT NULL,
    facebook_user_id VARCHAR(255) NOT NULL,
    facebook_user_name VARCHAR(500),
    source_type VARCHAR(20) NOT NULL CHECK (source_type IN ('comment', 'message')),
    post_id VARCHAR(255),
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_message_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'paused', 'closed')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_conversations_page_id ON conversations(page_id);
CREATE INDEX idx_conversations_facebook_user_id ON conversations(facebook_user_id);
CREATE INDEX idx_conversations_source_type ON conversations(source_type);
CREATE INDEX idx_conversations_status ON conversations(status);
CREATE INDEX idx_conversations_last_message_at ON conversations(last_message_at DESC);

-- Create unique constraint to prevent duplicate conversations
CREATE UNIQUE INDEX idx_conversations_unique_user_post ON conversations(page_id, facebook_user_id, COALESCE(post_id, ''));

-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_conversations_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_conversations_updated_at
    BEFORE UPDATE ON conversations
    FOR EACH ROW
    EXECUTE FUNCTION update_conversations_updated_at();
