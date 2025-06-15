-- Add missing columns to facebook_pages table
ALTER TABLE facebook_pages 
ADD COLUMN IF NOT EXISTS facebook_page_id VARCHAR(255),
ADD COLUMN IF NOT EXISTS name VARCHAR(255),
ADD COLUMN IF NOT EXISTS expires_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS facebook_user_id VARCHAR(255),
ADD COLUMN IF NOT EXISTS connected BOOLEAN DEFAULT TRUE,
ADD COLUMN IF NOT EXISTS enabled BOOLEAN DEFAULT TRUE,
ADD COLUMN IF NOT EXISTS webhook_verified BOOLEAN DEFAULT FALSE;

-- Update existing columns if they exist with different names
UPDATE facebook_pages SET 
  facebook_page_id = page_id,
  name = page_name,
  enabled = is_active
WHERE facebook_page_id IS NULL OR name IS NULL;

-- Create facebook_webhook_subscriptions table
CREATE TABLE IF NOT EXISTS facebook_webhook_subscriptions (
    id SERIAL PRIMARY KEY,
    facebook_page_id VARCHAR(255) NOT NULL,
    subscription_id VARCHAR(255),
    webhook_url TEXT NOT NULL,
    verify_token VARCHAR(255) NOT NULL,
    fields TEXT[] DEFAULT '{"comments", "mentions", "messages"}',
    active BOOLEAN DEFAULT TRUE,
    last_verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_facebook_pages_user_id ON facebook_pages(facebook_user_id);
CREATE INDEX IF NOT EXISTS idx_facebook_pages_facebook_page_id ON facebook_pages(facebook_page_id);
CREATE INDEX IF NOT EXISTS idx_webhook_subscriptions_page_id ON facebook_webhook_subscriptions(facebook_page_id);

-- Add foreign key constraint for webhook subscriptions
ALTER TABLE facebook_webhook_subscriptions 
ADD CONSTRAINT fk_webhook_subscriptions_page_id 
FOREIGN KEY (facebook_page_id) REFERENCES facebook_pages(facebook_page_id) ON DELETE CASCADE;

-- Add foreign key constraint for facebook_pages to facebook_sessions
ALTER TABLE facebook_pages 
ADD CONSTRAINT fk_facebook_pages_user_id 
FOREIGN KEY (facebook_user_id) REFERENCES facebook_sessions(facebook_user_id) ON DELETE CASCADE;

-- Create trigger for webhook subscriptions updated_at
CREATE TRIGGER update_facebook_webhook_subscriptions_updated_at 
BEFORE UPDATE ON facebook_webhook_subscriptions 
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
