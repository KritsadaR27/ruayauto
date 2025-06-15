-- Migration: Facebook Integration Tables
-- Created: June 15, 2025

-- Facebook user sessions table
CREATE TABLE IF NOT EXISTS facebook_sessions (
    id SERIAL PRIMARY KEY,
    facebook_user_id VARCHAR(255) NOT NULL UNIQUE,
    access_token TEXT NOT NULL,
    expires_at TIMESTAMP,
    refresh_token TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Facebook pages table
CREATE TABLE IF NOT EXISTS facebook_pages (
    id SERIAL PRIMARY KEY,
    facebook_page_id VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    access_token TEXT NOT NULL,
    expires_at TIMESTAMP,
    facebook_user_id VARCHAR(255) NOT NULL REFERENCES facebook_sessions(facebook_user_id) ON DELETE CASCADE,
    connected BOOLEAN DEFAULT TRUE,
    enabled BOOLEAN DEFAULT TRUE,
    webhook_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for facebook_pages table
CREATE INDEX IF NOT EXISTS idx_facebook_pages_user_id ON facebook_pages(facebook_user_id);
CREATE INDEX IF NOT EXISTS idx_facebook_pages_page_id ON facebook_pages(facebook_page_id);

-- Facebook webhook subscriptions table
CREATE TABLE IF NOT EXISTS facebook_webhook_subscriptions (
    id SERIAL PRIMARY KEY,
    facebook_page_id VARCHAR(255) NOT NULL REFERENCES facebook_pages(facebook_page_id) ON DELETE CASCADE,
    subscription_id VARCHAR(255),
    webhook_url TEXT NOT NULL,
    verify_token VARCHAR(255) NOT NULL,
    fields TEXT[] DEFAULT '{"comments", "mentions", "messages"}',
    active BOOLEAN DEFAULT TRUE,
    last_verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for facebook_webhook_subscriptions table  
CREATE INDEX IF NOT EXISTS idx_webhook_subscriptions_page_id ON facebook_webhook_subscriptions(facebook_page_id);

-- Update triggers for updated_at columns
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_facebook_sessions_updated_at BEFORE UPDATE ON facebook_sessions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_facebook_pages_updated_at BEFORE UPDATE ON facebook_pages FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_facebook_webhook_subscriptions_updated_at BEFORE UPDATE ON facebook_webhook_subscriptions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
