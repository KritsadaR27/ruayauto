-- Create chatbot_mvp schema for MVP development
CREATE SCHEMA IF NOT EXISTS chatbot_mvp;

-- Set search path to use the new schema
SET search_path TO chatbot_mvp, public;

-- Conversations table - stores Facebook user conversations
CREATE TABLE chatbot_mvp.conversations (
    id SERIAL PRIMARY KEY,
    facebook_user_id VARCHAR(255) UNIQUE NOT NULL,
    user_name VARCHAR(255),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'paused', 'closed')),
    last_message_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create index for faster lookups
CREATE INDEX idx_conversations_facebook_user_id ON chatbot_mvp.conversations(facebook_user_id);
CREATE INDEX idx_conversations_status ON chatbot_mvp.conversations(status);
CREATE INDEX idx_conversations_last_message ON chatbot_mvp.conversations(last_message_at DESC);

-- Messages table - stores all messages in conversations
CREATE TABLE chatbot_mvp.messages (
    id SERIAL PRIMARY KEY,
    conversation_id INTEGER REFERENCES chatbot_mvp.conversations(id) ON DELETE CASCADE,
    facebook_message_id VARCHAR(255),
    sender_type VARCHAR(10) NOT NULL CHECK (sender_type IN ('user', 'bot')),
    message_type VARCHAR(20) DEFAULT 'text' CHECK (message_type IN ('text', 'image', 'quick_reply', 'template', 'audio', 'video', 'file')),
    message_text TEXT,
    message_data JSONB,
    matched_keyword VARCHAR(255),
    response_template VARCHAR(255),
    processing_time_ms INTEGER,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for messages
CREATE INDEX idx_messages_conversation_id ON chatbot_mvp.messages(conversation_id);
CREATE INDEX idx_messages_created_at ON chatbot_mvp.messages(created_at DESC);
CREATE INDEX idx_messages_sender_type ON chatbot_mvp.messages(sender_type);
CREATE INDEX idx_messages_facebook_message_id ON chatbot_mvp.messages(facebook_message_id);

-- Keywords table - enhanced keyword matching
CREATE TABLE chatbot_mvp.keywords (
    id SERIAL PRIMARY KEY,
    keyword TEXT NOT NULL,
    response_type VARCHAR(20) DEFAULT 'text' CHECK (response_type IN ('text', 'image', 'quick_reply', 'template')),
    response_content TEXT NOT NULL,
    response_data JSONB,
    priority INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    match_type VARCHAR(20) DEFAULT 'exact' CHECK (match_type IN ('exact', 'partial', 'regex', 'fuzzy')),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for keywords
CREATE INDEX idx_keywords_keyword ON chatbot_mvp.keywords(keyword);
CREATE INDEX idx_keywords_is_active ON chatbot_mvp.keywords(is_active);
CREATE INDEX idx_keywords_priority ON chatbot_mvp.keywords(priority DESC);

-- Response templates table - reusable response templates
CREATE TABLE chatbot_mvp.response_templates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    template_type VARCHAR(20) NOT NULL CHECK (template_type IN ('text', 'image', 'quick_reply', 'generic', 'button', 'carousel')),
    template_data JSONB NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for response templates
CREATE INDEX idx_response_templates_name ON chatbot_mvp.response_templates(name);
CREATE INDEX idx_response_templates_is_active ON chatbot_mvp.response_templates(is_active);
CREATE INDEX idx_response_templates_template_type ON chatbot_mvp.response_templates(template_type);

-- Conversation stats table - analytics data
CREATE TABLE chatbot_mvp.conversation_stats (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL,
    total_conversations INTEGER DEFAULT 0,
    total_messages INTEGER DEFAULT 0,
    bot_responses INTEGER DEFAULT 0,
    avg_response_time_ms DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(date)
);

-- Create index for stats
CREATE INDEX idx_conversation_stats_date ON chatbot_mvp.conversation_stats(date DESC);

-- Insert some sample keywords for testing
INSERT INTO chatbot_mvp.keywords (keyword, response_type, response_content, priority) VALUES
('สวัสดี', 'text', 'สวัสดีค่ะ! ยินดีต้อนรับครับ 😊', 10),
('ขอบคุณ', 'text', 'ยินดีค่ะ! มีอะไรให้ช่วยอีกไหมคะ', 5),
('ราคา', 'text', 'สำหรับข้อมูลราคา กรุณาแจ้งสินค้าที่สนใจค่ะ', 5),
('โปร', 'text', 'โปรโมชั่นประจำเดือนนี้: ซื้อ 2 แถม 1 ค่ะ! 🎉', 8),
('ส่ง', 'text', 'เราส่งทั่วประเทศไทยค่ะ ค่าส่งเริ่มต้น 50 บาท 📦', 6);

-- Insert sample response templates
INSERT INTO chatbot_mvp.response_templates (name, template_type, template_data, description) VALUES
('welcome_message', 'text', '{"text": "ยินดีต้อนรับค่ะ! มีอะไรให้ช่วยไหมคะ"}', 'Welcome message for new users'),
('quick_reply_menu', 'quick_reply', '{"text": "เลือกหัวข้อที่สนใจค่ะ", "quick_replies": [{"title": "ดูสินค้า", "payload": "VIEW_PRODUCTS"}, {"title": "สอบถาม", "payload": "ASK_QUESTION"}, {"title": "ติดต่อแอดมิน", "payload": "CONTACT_ADMIN"}]}', 'Main menu with quick reply options');

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION chatbot_mvp.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for auto-updating updated_at
CREATE TRIGGER update_conversations_updated_at BEFORE UPDATE ON chatbot_mvp.conversations FOR EACH ROW EXECUTE FUNCTION chatbot_mvp.update_updated_at_column();
CREATE TRIGGER update_keywords_updated_at BEFORE UPDATE ON chatbot_mvp.keywords FOR EACH ROW EXECUTE FUNCTION chatbot_mvp.update_updated_at_column();
