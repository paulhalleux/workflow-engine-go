CREATE TABLE IF NOT EXISTS type_schemas (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name TEXT NOT NULL,
    version TEXT NOT NULL,
    schema JSONB NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_type_schemas_name_version ON type_schemas (name, version);

ALTER TABLE workflow_definitions ADD COLUMN IF NOT EXISTS input_parameters JSONB;
ALTER TABLE workflow_definitions ADD COLUMN IF NOT EXISTS output_parameters JSONB;