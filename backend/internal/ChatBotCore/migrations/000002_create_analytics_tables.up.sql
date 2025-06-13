-- Create analytics table for tracking message interactions
CREATE TABLE IF NOT EXISTS message_analytics (
    id SERIAL PRIMARY KEY,
    facebook_message_id VARCHAR(255),
    sender_id VARCHAR(255) NOT NULL,
    recipient_id VARCHAR(255) NOT NULL,
    message_text TEXT,
    matched_keyword VARCHAR(255),
    response_sent TEXT,
    response_time_ms INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for analytics queries
CREATE INDEX IF NOT EXISTS idx_analytics_sender ON message_analytics(sender_id);
CREATE INDEX IF NOT EXISTS idx_analytics_recipient ON message_analytics(recipient_id);
CREATE INDEX IF NOT EXISTS idx_analytics_keyword ON message_analytics(matched_keyword);
CREATE INDEX IF NOT EXISTS idx_analytics_created_at ON message_analytics(created_at);

-- Create webhook_logs table for debugging
CREATE TABLE IF NOT EXISTS webhook_logs (
    id SERIAL PRIMARY KEY,
    webhook_type VARCHAR(50) NOT NULL,
    request_body JSONB,
    response_status INTEGER,
    response_body TEXT,
    processing_time_ms INTEGER,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for webhook logs
CREATE INDEX IF NOT EXISTS idx_webhook_logs_type ON webhook_logs(webhook_type);
CREATE INDEX IF NOT EXISTS idx_webhook_logs_created_at ON webhook_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_webhook_logs_status ON webhook_logs(response_status);
