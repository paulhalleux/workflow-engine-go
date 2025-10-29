DROP TABLE IF EXISTS type_schemas;
ALTER TABLE workflow_definitions DROP COLUMN IF EXISTS input_parameters;
ALTER TABLE workflow_definitions DROP COLUMN IF EXISTS output_parameters;
DROP INDEX IF EXISTS idx_type_schemas_name_version;