-- Create response_templates table for storing reusable response templates
CREATE TABLE IF NOT EXISTS response_templates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    response_type VARCHAR(20) NOT NULL CHECK (response_type IN ('text', 'image', 'button', 'carousel', 'quick_reply')),
    content TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_response_templates_name ON response_templates(name);
CREATE INDEX idx_response_templates_response_type ON response_templates(response_type);
CREATE INDEX idx_response_templates_active ON response_templates(active);

-- Create index on metadata for JSON queries
CREATE INDEX idx_response_templates_metadata_gin ON response_templates USING gin(metadata);

-- Create unique constraint on template name
CREATE UNIQUE INDEX idx_response_templates_name_unique ON response_templates(name) WHERE active = true;

-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_response_templates_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_response_templates_updated_at
    BEFORE UPDATE ON response_templates
    FOR EACH ROW
    EXECUTE FUNCTION update_response_templates_updated_at();
