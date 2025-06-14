-- Migration: Create Analytics and Rate Limiting Tables
-- Description: Add comprehensive analytics and rate limiting capabilities
-- Phase 2: Analytics and Performance Monitoring

-- 1. Create rule_analytics table for detailed usage tracking
CREATE TABLE IF NOT EXISTS rule_analytics (
    id SERIAL PRIMARY KEY,
    rule_id INTEGER REFERENCES rules(id) ON DELETE CASCADE,
    page_id INTEGER REFERENCES pages(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    
    -- Event details
    event_type VARCHAR(50) NOT NULL, -- 'comment_matched', 'response_sent', 'rule_triggered', 'condition_failed'
    
    -- Comment/Response details
    comment_id VARCHAR(255),
    comment_text TEXT,
    matched_keywords JSONB DEFAULT '[]',
    response_sent TEXT,
    response_type VARCHAR(50) DEFAULT 'text',
    
    -- Performance metrics
    processing_time_ms INTEGER,
    confidence_score DECIMAL(5,4) DEFAULT 0.0,
    
    -- Context data
    comment_author_id VARCHAR(255),
    comment_author_name VARCHAR(255),
    post_id VARCHAR(255),
    
    -- Success/failure tracking
    is_successful BOOLEAN DEFAULT TRUE,
    error_message TEXT,
    
    -- Metadata
    user_agent TEXT,
    ip_address INET,
    metadata JSONB DEFAULT '{}',
    
    created_at TIMESTAMP DEFAULT NOW()
);

-- 2. Create rate_limits table for tracking user interaction limits
CREATE TABLE IF NOT EXISTS rate_limits (
    id SERIAL PRIMARY KEY,
    rule_id INTEGER REFERENCES rules(id) ON DELETE CASCADE,
    page_id INTEGER REFERENCES pages(id) ON DELETE CASCADE,
    
    -- User identification
    facebook_user_id VARCHAR(255) NOT NULL,
    facebook_user_name VARCHAR(255),
    
    -- Rate limit tracking
    interaction_count INTEGER DEFAULT 1,
    first_interaction_at TIMESTAMP DEFAULT NOW(),
    last_interaction_at TIMESTAMP DEFAULT NOW(),
    
    -- Rate limit settings (copied from rule for historical tracking)
    max_interactions INTEGER DEFAULT 3,
    window_hours INTEGER DEFAULT 24,
    
    -- Status
    is_blocked BOOLEAN DEFAULT FALSE,
    blocked_until TIMESTAMP,
    blocked_reason TEXT,
    
    -- Metadata
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    -- Composite unique constraint
    UNIQUE(rule_id, page_id, facebook_user_id)
);

-- 3. Create analytics_summary table for aggregated metrics
CREATE TABLE IF NOT EXISTS analytics_summary (
    id SERIAL PRIMARY KEY,
    
    -- Grouping dimensions
    date DATE NOT NULL,
    rule_id INTEGER REFERENCES rules(id) ON DELETE CASCADE,
    page_id INTEGER REFERENCES pages(id) ON DELETE CASCADE,
    
    -- Metrics
    total_comments INTEGER DEFAULT 0,
    total_matches INTEGER DEFAULT 0,
    total_responses INTEGER DEFAULT 0,
    total_errors INTEGER DEFAULT 0,
    
    -- Performance metrics
    avg_processing_time_ms DECIMAL(8,2) DEFAULT 0.0,
    avg_confidence_score DECIMAL(5,4) DEFAULT 0.0,
    
    -- Response type breakdown
    text_responses INTEGER DEFAULT 0,
    image_responses INTEGER DEFAULT 0,
    mixed_responses INTEGER DEFAULT 0,
    
    -- Rate limiting stats
    rate_limited_users INTEGER DEFAULT 0,
    total_blocked_interactions INTEGER DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    -- Unique constraint for daily summaries
    UNIQUE(date, rule_id, page_id)
);

-- 4. Create performance_logs table for system monitoring
CREATE TABLE IF NOT EXISTS performance_logs (
    id SERIAL PRIMARY KEY,
    
    -- System metrics
    endpoint VARCHAR(255),
    method VARCHAR(10),
    response_time_ms INTEGER,
    status_code INTEGER,
    
    -- Resource usage
    memory_usage_mb DECIMAL(10,2),
    cpu_usage_percent DECIMAL(5,2),
    
    -- Request details
    request_size_bytes INTEGER,
    response_size_bytes INTEGER,
    
    -- User context
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    page_id INTEGER REFERENCES pages(id) ON DELETE SET NULL,
    
    -- Error tracking
    error_message TEXT,
    stack_trace TEXT,
    
    -- Metadata
    ip_address INET,
    user_agent TEXT,
    metadata JSONB DEFAULT '{}',
    
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_rule_analytics_rule_id ON rule_analytics(rule_id);
CREATE INDEX IF NOT EXISTS idx_rule_analytics_page_id ON rule_analytics(page_id);
CREATE INDEX IF NOT EXISTS idx_rule_analytics_event_type ON rule_analytics(event_type);
CREATE INDEX IF NOT EXISTS idx_rule_analytics_created_at ON rule_analytics(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_rule_analytics_comment_id ON rule_analytics(comment_id);

CREATE INDEX IF NOT EXISTS idx_rate_limits_rule_page_user ON rate_limits(rule_id, page_id, facebook_user_id);
CREATE INDEX IF NOT EXISTS idx_rate_limits_facebook_user_id ON rate_limits(facebook_user_id);
CREATE INDEX IF NOT EXISTS idx_rate_limits_is_blocked ON rate_limits(is_blocked);
CREATE INDEX IF NOT EXISTS idx_rate_limits_last_interaction ON rate_limits(last_interaction_at);

CREATE INDEX IF NOT EXISTS idx_analytics_summary_date ON analytics_summary(date DESC);
CREATE INDEX IF NOT EXISTS idx_analytics_summary_rule_id ON analytics_summary(rule_id);
CREATE INDEX IF NOT EXISTS idx_analytics_summary_page_id ON analytics_summary(page_id);

CREATE INDEX IF NOT EXISTS idx_performance_logs_endpoint ON performance_logs(endpoint);
CREATE INDEX IF NOT EXISTS idx_performance_logs_created_at ON performance_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_performance_logs_status_code ON performance_logs(status_code);
CREATE INDEX IF NOT EXISTS idx_performance_logs_response_time ON performance_logs(response_time_ms);

-- Add trigger for rate_limits updated_at
CREATE TRIGGER update_rate_limits_updated_at 
    BEFORE UPDATE ON rate_limits 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_analytics_summary_updated_at 
    BEFORE UPDATE ON analytics_summary 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create view for real-time analytics dashboard
CREATE OR REPLACE VIEW analytics_dashboard AS
SELECT 
    r.id as rule_id,
    r.name as rule_name,
    p.name as page_name,
    p.platform,
    
    -- Today's metrics
    COALESCE(today.total_comments, 0) as today_comments,
    COALESCE(today.total_matches, 0) as today_matches,
    COALESCE(today.total_responses, 0) as today_responses,
    COALESCE(today.total_errors, 0) as today_errors,
    
    -- Overall metrics
    r.usage_count as total_usage,
    r.last_used_at,
    
    -- Success rate
    CASE 
        WHEN COALESCE(today.total_comments, 0) > 0 
        THEN ROUND((COALESCE(today.total_matches, 0)::DECIMAL / today.total_comments) * 100, 2)
        ELSE 0
    END as today_success_rate,
    
    -- Performance metrics
    COALESCE(today.avg_processing_time_ms, 0) as avg_processing_time,
    COALESCE(today.avg_confidence_score, 0) as avg_confidence_score,
    
    -- Rate limiting stats
    COALESCE(today.rate_limited_users, 0) as rate_limited_users,
    
    r.is_active,
    r.priority

FROM rules r
LEFT JOIN pages p ON EXISTS (
    SELECT 1 FROM rule_pages rp 
    WHERE rp.rule_id = r.id AND rp.page_id = p.id AND rp.is_active = TRUE
)
LEFT JOIN analytics_summary today ON (
    today.rule_id = r.id 
    AND today.page_id = p.id 
    AND today.date = CURRENT_DATE
)
WHERE r.is_active = TRUE
ORDER BY r.priority DESC, r.usage_count DESC;

-- Create stored procedure for cleaning old analytics data
CREATE OR REPLACE FUNCTION cleanup_old_analytics(retention_days INTEGER DEFAULT 90)
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    -- Delete old rule_analytics records
    DELETE FROM rule_analytics 
    WHERE created_at < NOW() - INTERVAL '1 day' * retention_days;
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    
    -- Delete old performance_logs
    DELETE FROM performance_logs 
    WHERE created_at < NOW() - INTERVAL '1 day' * retention_days;
    
    -- Clean up old rate_limits for unblocked users
    DELETE FROM rate_limits 
    WHERE is_blocked = FALSE 
    AND last_interaction_at < NOW() - INTERVAL '1 day' * 30; -- Keep rate limits for 30 days
    
    -- Log cleanup activity
    INSERT INTO performance_logs (
        endpoint, method, response_time_ms, status_code,
        metadata, created_at
    ) VALUES (
        'cleanup_old_analytics', 'SYSTEM', 0, 200,
        json_build_object('deleted_records', deleted_count, 'retention_days', retention_days),
        NOW()
    );
    
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Insert sample analytics data for testing
INSERT INTO rule_analytics (
    rule_id, page_id, event_type, comment_text, 
    matched_keywords, response_sent, processing_time_ms, 
    confidence_score, is_successful, created_at
)
SELECT 
    r.id,
    p.id,
    'comment_matched',
    'Sample comment for testing rule: ' || r.name,
    '["' || k.keyword || '"]'::jsonb,
    rr.response_text,
    FLOOR(RANDOM() * 100 + 10)::INTEGER, -- 10-110ms
    ROUND((RANDOM() * 0.5 + 0.5)::DECIMAL, 4), -- 0.5-1.0 confidence
    TRUE,
    NOW() - INTERVAL '1 hour' * FLOOR(RANDOM() * 24) -- Random time in last 24 hours
FROM rules r
JOIN rule_pages rp ON r.id = rp.rule_id
JOIN pages p ON rp.page_id = p.id
JOIN rule_keywords rk ON r.id = rk.rule_id
JOIN keywords k ON rk.keyword_id = k.id
JOIN rule_responses rr ON r.id = rr.rule_id
WHERE r.is_active = TRUE
LIMIT 100; -- Add 100 sample analytics records

COMMENT ON TABLE rule_analytics IS 'Detailed analytics for rule usage and performance';
COMMENT ON TABLE rate_limits IS 'Rate limiting tracking per user per rule';
COMMENT ON TABLE analytics_summary IS 'Daily aggregated analytics summary';
COMMENT ON TABLE performance_logs IS 'System performance monitoring logs';
COMMENT ON VIEW analytics_dashboard IS 'Real-time analytics dashboard view';
