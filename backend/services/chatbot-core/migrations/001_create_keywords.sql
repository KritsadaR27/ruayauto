-- Create keywords table for auto-reply functionality
CREATE TABLE IF NOT EXISTS keywords (
    id SERIAL PRIMARY KEY,
    keyword VARCHAR(255) NOT NULL,
    response TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    priority INTEGER DEFAULT 0,
    match_type VARCHAR(50) DEFAULT 'contains' CHECK (match_type IN ('exact', 'contains', 'starts_with', 'ends_with')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for better performance
CREATE INDEX IF NOT EXISTS idx_keywords_active ON keywords(is_active);
CREATE INDEX IF NOT EXISTS idx_keywords_priority ON keywords(priority DESC);

-- Insert some sample keywords for testing
INSERT INTO keywords (keyword, response, is_active, priority, match_type) VALUES
('สวัสดี', 'สวัสดีครับ! ยินดีให้บริการครับ', true, 100, 'contains'),
('ขอบคุณ', 'ยินดีครับ! หากมีคำถามเพิ่มเติม สามารถสอบถามได้เสมอครับ', true, 90, 'contains'),
('ราคา', 'สำหรับข้อมูลราคา กรุณาติดต่อทีมขายของเราโดยตรงครับ', true, 80, 'contains'),
('โปรโมชั่น', 'ขณะนี้มีโปรโมชั่นพิเศษมากมาย ติดตามได้ที่เพจของเราครับ', true, 85, 'contains'),
('ขอใบเสนอราคา', 'ขอบคุณที่สนใจ เราจะให้ทีมงานติดต่อกลับไปเพื่อส่งใบเสนอราคาให้ครับ', true, 95, 'contains');
