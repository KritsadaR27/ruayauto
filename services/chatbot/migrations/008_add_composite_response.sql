-- Migration: Add composite response support
-- Allow text + media in single response

ALTER TABLE rule_responses 
ADD COLUMN has_media BOOLEAN DEFAULT FALSE,
ADD COLUMN media_description TEXT;

-- Update existing records
UPDATE rule_responses 
SET has_media = CASE 
    WHEN media_url IS NOT NULL AND media_url != '' THEN TRUE 
    ELSE FALSE 
END;

-- Add comment for clarity
COMMENT ON COLUMN rule_responses.response_type IS 'Primary type: text, image, video, or composite';
COMMENT ON COLUMN rule_responses.has_media IS 'Whether this response includes media attachment';
COMMENT ON COLUMN rule_responses.media_description IS 'Alt text or description for media';
