-- Create keyword_responses junction table for many-to-many relationship
CREATE TABLE IF NOT EXISTS keyword_responses (
    keyword_id INTEGER NOT NULL REFERENCES keywords(id) ON DELETE CASCADE,
    response_template_id INTEGER NOT NULL REFERENCES response_templates(id) ON DELETE CASCADE,
    order_index INTEGER DEFAULT 0,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (keyword_id, response_template_id)
);

-- Create indexes for better query performance
CREATE INDEX idx_keyword_responses_keyword_id ON keyword_responses(keyword_id);
CREATE INDEX idx_keyword_responses_response_template_id ON keyword_responses(response_template_id);
CREATE INDEX idx_keyword_responses_order_index ON keyword_responses(order_index);
CREATE INDEX idx_keyword_responses_active ON keyword_responses(active);

-- Create composite index for active responses ordered by index
CREATE INDEX idx_keyword_responses_active_order ON keyword_responses(keyword_id, active, order_index) WHERE active = true;

-- Add foreign key constraints with proper naming
ALTER TABLE keyword_responses 
ADD CONSTRAINT fk_keyword_responses_keyword_id 
FOREIGN KEY (keyword_id) REFERENCES keywords(id) ON DELETE CASCADE;

ALTER TABLE keyword_responses 
ADD CONSTRAINT fk_keyword_responses_response_template_id 
FOREIGN KEY (response_template_id) REFERENCES response_templates(id) ON DELETE CASCADE;
