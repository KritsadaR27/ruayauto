-- Update existing keywords table to match new requirements
ALTER TABLE keywords 
ADD COLUMN IF NOT EXISTS match_type VARCHAR(20) DEFAULT 'exact' CHECK (match_type IN ('exact', 'contains')),
ADD COLUMN IF NOT EXISTS priority INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS active BOOLEAN DEFAULT true;

-- Create index for better query performance
CREATE INDEX IF NOT EXISTS idx_keywords_match_type ON keywords(match_type);
CREATE INDEX IF NOT EXISTS idx_keywords_priority ON keywords(priority DESC);
CREATE INDEX IF NOT EXISTS idx_keywords_active ON keywords(active);

-- Create composite index for active keywords ordered by priority
CREATE INDEX IF NOT EXISTS idx_keywords_active_priority ON keywords(active, priority DESC) WHERE active = true;

-- Create text search index for keyword matching
CREATE INDEX IF NOT EXISTS idx_keywords_keyword_text ON keywords USING gin(to_tsvector('english', keyword));

-- Update existing records to have default values
UPDATE keywords 
SET 
    match_type = 'exact',
    priority = 0,
    active = COALESCE(is_active, true)
WHERE match_type IS NULL OR priority IS NULL OR active IS NULL;
