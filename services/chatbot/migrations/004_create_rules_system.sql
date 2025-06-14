-- Migration 004: Enhanced Rules System with Multi-Keywords and Multi-Responses
-- Building on existing keywords, messages, facebook_pages tables

-- Add new columns to existing keywords table for enhanced features (skip if exists)
ALTER TABLE keywords ADD COLUMN IF NOT EXISTS rule_name VARCHAR(255) DEFAULT '';
ALTER TABLE keywords ADD COLUMN IF NOT EXISTS hide_after_reply BOOLEAN DEFAULT FALSE;
ALTER TABLE keywords ADD COLUMN IF NOT EXISTS send_to_inbox BOOLEAN DEFAULT FALSE;
ALTER TABLE keywords ADD COLUMN IF NOT EXISTS inbox_message TEXT DEFAULT '';
ALTER TABLE keywords ADD COLUMN IF NOT EXISTS inbox_image VARCHAR(500) DEFAULT '';
ALTER TABLE keywords ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT TRUE;
ALTER TABLE keywords ADD COLUMN IF NOT EXISTS priority INTEGER DEFAULT 1;
-- Skip created_by for now since users table may not exist
-- ALTER TABLE keywords ADD COLUMN created_by INTEGER REFERENCES users(id);
ALTER TABLE keywords ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Create keyword_responses table for multiple responses per rule
CREATE TABLE keyword_responses (
    id SERIAL PRIMARY KEY,
    keyword_id INTEGER NOT NULL REFERENCES keywords(id) ON DELETE CASCADE,
    response_text TEXT NOT NULL,
    response_type VARCHAR(50) DEFAULT 'text', -- text, image, video, link
    media_url VARCHAR(500),
    weight INTEGER DEFAULT 1, -- For weighted random selection
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create keyword_pages table for page-specific rules  
CREATE TABLE keyword_pages (
    id SERIAL PRIMARY KEY,
    keyword_id INTEGER NOT NULL REFERENCES keywords(id) ON DELETE CASCADE,
    page_id INTEGER NOT NULL REFERENCES facebook_pages(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(keyword_id, page_id)
);

-- Create fallback_responses table for when no keywords match
CREATE TABLE fallback_responses (
    id SERIAL PRIMARY KEY,
    page_id INTEGER REFERENCES facebook_pages(id) ON DELETE CASCADE,
    response_text TEXT NOT NULL,
    response_type VARCHAR(50) DEFAULT 'text',
    media_url VARCHAR(500),
    weight INTEGER DEFAULT 1,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create rate_limits table to prevent spam
CREATE TABLE rate_limits (
    id SERIAL PRIMARY KEY,
    page_id INTEGER NOT NULL REFERENCES facebook_pages(id) ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL, -- Facebook user ID
    keyword_id INTEGER REFERENCES keywords(id) ON DELETE CASCADE,
    last_response_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    response_count INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(page_id, user_id, keyword_id)
);

-- Create scheduled_messages table for delayed responses
CREATE TABLE scheduled_messages (
    id SERIAL PRIMARY KEY,
    page_id INTEGER NOT NULL REFERENCES facebook_pages(id) ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL,
    message_text TEXT NOT NULL,
    media_url VARCHAR(500),
    scheduled_at TIMESTAMP NOT NULL,
    status VARCHAR(50) DEFAULT 'pending', -- pending, sent, failed
    attempts INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create incoming_messages table for analytics and inbox features
CREATE TABLE incoming_messages (
    id SERIAL PRIMARY KEY,
    page_id INTEGER NOT NULL REFERENCES facebook_pages(id) ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL,
    message_text TEXT NOT NULL,
    comment_id VARCHAR(255),
    post_id VARCHAR(255),
    matched_keyword_id INTEGER REFERENCES keywords(id),
    response_sent BOOLEAN DEFAULT FALSE,
    hidden_comment BOOLEAN DEFAULT FALSE,
    sent_to_inbox BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create rule_analytics table for performance tracking
CREATE TABLE rule_analytics (
    id SERIAL PRIMARY KEY,
    keyword_id INTEGER NOT NULL REFERENCES keywords(id) ON DELETE CASCADE,
    page_id INTEGER NOT NULL REFERENCES facebook_pages(id) ON DELETE CASCADE,
    trigger_count INTEGER DEFAULT 0,
    response_count INTEGER DEFAULT 0,
    last_triggered_at TIMESTAMP,
    date DATE DEFAULT CURRENT_DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(keyword_id, page_id, date)
);

-- Create view for easy data retrieval
CREATE OR REPLACE VIEW keywords_with_responses AS
SELECT 
    k.id,
    k.keyword,
    k.response as default_message,
    k.rule_name,
    k.hide_after_reply,
    k.send_to_inbox,
    k.inbox_message,
    k.inbox_image,
    k.is_active,
    k.priority,
    NULL::integer as created_by,
    k.created_at,
    k.updated_at,
    COALESCE(
        json_agg(
            json_build_object(
                'id', kr.id,
                'response_text', kr.response_text,
                'response_type', kr.response_type,
                'media_url', kr.media_url,
                'weight', kr.weight
            ) ORDER BY kr.weight DESC
        ) FILTER (WHERE kr.id IS NOT NULL),
        '[]'::json
    ) as responses,
    COALESCE(
        array_agg(DISTINCT fp.page_id) FILTER (WHERE fp.page_id IS NOT NULL),
        '{}'::integer[]
    ) as page_ids
FROM keywords k
LEFT JOIN keyword_responses kr ON k.id = kr.keyword_id AND kr.is_active = TRUE
LEFT JOIN keyword_pages kp ON k.id = kp.keyword_id
LEFT JOIN facebook_pages fp ON kp.page_id = fp.id
WHERE k.is_active = TRUE
GROUP BY k.id, k.keyword, k.message, k.rule_name, k.hide_after_reply, 
         k.send_to_inbox, k.inbox_message, k.inbox_image, k.is_active, 
         k.priority, k.created_by, k.created_at, k.updated_at;

-- Create indexes for performance (skip if exists)
CREATE INDEX IF NOT EXISTS idx_keywords_active ON keywords(is_active);
CREATE INDEX IF NOT EXISTS idx_keywords_priority ON keywords(priority DESC);
CREATE INDEX IF NOT EXISTS idx_keyword_responses_keyword_id ON keyword_responses(keyword_id);
CREATE INDEX IF NOT EXISTS idx_keyword_responses_active ON keyword_responses(is_active);
CREATE INDEX IF NOT EXISTS idx_keyword_pages_keyword_id ON keyword_pages(keyword_id);
CREATE INDEX IF NOT EXISTS idx_keyword_pages_page_id ON keyword_pages(page_id);
CREATE INDEX IF NOT EXISTS idx_rate_limits_page_user ON rate_limits(page_id, user_id);
CREATE INDEX IF NOT EXISTS idx_rate_limits_last_response ON rate_limits(last_response_at);
CREATE INDEX IF NOT EXISTS idx_incoming_messages_page_id ON incoming_messages(page_id);
CREATE INDEX IF NOT EXISTS idx_incoming_messages_created_at ON incoming_messages(created_at);
CREATE INDEX IF NOT EXISTS idx_rule_analytics_keyword_page ON rule_analytics(keyword_id, page_id);
CREATE INDEX IF NOT EXISTS idx_rule_analytics_date ON rule_analytics(date);
CREATE INDEX IF NOT EXISTS idx_scheduled_messages_status ON scheduled_messages(status);
CREATE INDEX IF NOT EXISTS idx_scheduled_messages_scheduled_at ON scheduled_messages(scheduled_at);

-- Insert some default fallback responses
INSERT INTO fallback_responses (response_text, weight) VALUES
('ขอบคุณที่ติดต่อเรา เราจะตอบกลับเร็วๆ นี้ครับ', 3),
('สวัสดีครับ ยินดีให้บริการ', 2),
('ขอบคุณสำหรับความสนใจ เราจะติดต่อกลับเร็วๆ นี้', 2),
('หากต้องการความช่วยเหลือเพิ่มเติม กรุณาส่งข้อความส่วนตัวมาครับ', 1);

-- Migration complete message (skip if messages table doesn't have message column)
-- INSERT INTO messages (message, created_at) VALUES 
-- ('Migration 004: Enhanced Rules System completed successfully', CURRENT_TIMESTAMP);

-- Log migration completion to a simple way
DO $$ 
BEGIN 
    RAISE NOTICE 'Migration 004: Enhanced Rules System completed successfully';
END $$;

COMMENT ON TABLE keyword_responses IS 'Multiple responses for each keyword rule';
COMMENT ON TABLE keyword_pages IS 'Page-specific keyword rules';
COMMENT ON TABLE fallback_responses IS 'Default responses when no keywords match';
COMMENT ON TABLE rate_limits IS 'Rate limiting to prevent spam responses';
COMMENT ON TABLE scheduled_messages IS 'Messages scheduled for future delivery';
COMMENT ON TABLE incoming_messages IS 'Log of all incoming messages for analytics';
COMMENT ON TABLE rule_analytics IS 'Performance analytics for keyword rules';