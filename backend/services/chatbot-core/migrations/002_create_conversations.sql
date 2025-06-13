-- Create conversations table
CREATE TABLE IF NOT EXISTS conversations (
    id SERIAL PRIMARY KEY,
    page_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    comment_id VARCHAR(255),
    post_id VARCHAR(255),
    status VARCHAR(50) DEFAULT 'active',
    last_message TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_conversations_page_user ON conversations(page_id, user_id);
CREATE INDEX IF NOT EXISTS idx_conversations_status ON conversations(status);
CREATE INDEX IF NOT EXISTS idx_conversations_last_message ON conversations(last_message);
