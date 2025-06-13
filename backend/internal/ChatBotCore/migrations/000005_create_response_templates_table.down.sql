-- Drop response_templates table and related objects
DROP TRIGGER IF EXISTS trigger_response_templates_updated_at ON response_templates;
DROP FUNCTION IF EXISTS update_response_templates_updated_at();
DROP INDEX IF EXISTS idx_response_templates_name_unique;
DROP INDEX IF EXISTS idx_response_templates_metadata_gin;
DROP INDEX IF EXISTS idx_response_templates_active;
DROP INDEX IF EXISTS idx_response_templates_response_type;
DROP INDEX IF EXISTS idx_response_templates_name;
DROP TABLE IF EXISTS response_templates;
