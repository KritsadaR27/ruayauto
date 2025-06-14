-- Create messages table for storing all incoming messages for AI analysis
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    platform VARCHAR(50) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    page_id VARCHAR(255),
    content TEXT NOT NULL,
    comment_id VARCHAR(255),
    post_id VARCHAR(255),
    message_type VARCHAR(50) DEFAULT 'comment', -- comment, inbox, post
    processed BOOLEAN DEFAULT FALSE,
    response TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create facebook_pages table for managing connected pages
CREATE TABLE IF NOT EXISTS facebook_pages (
    id SERIAL PRIMARY KEY,
    page_id VARCHAR(255) UNIQUE NOT NULL,
    page_name VARCHAR(255) NOT NULL,
    access_token TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_messages_platform ON messages(platform);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);
CREATE INDEX IF NOT EXISTS idx_messages_processed ON messages(processed);
CREATE INDEX IF NOT EXISTS idx_facebook_pages_active ON facebook_pages(is_active);

-- Insert some sample keywords for testing (using rare test words to avoid false positives)
INSERT INTO keywords (keyword, response, is_active, priority, match_type) VALUES
('xyztest123', 'Test greeting response activated', true, 10, 'contains'),
('qwertyuiop987', 'Alternative greeting test response', true, 10, 'contains'),
('abcdefg789', 'Test pricing response activated', true, 8, 'contains'),
('mnbvcxz321', 'Alternative pricing test response', true, 8, 'contains'),
('poiuytrewq654', 'Thank you test response activated', true, 5, 'contains')
ON CONFLICT DO NOTHING;
