CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS workflow_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    version VARCHAR(50) NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    input_parameters JSONB,
    output_parameters JSONB,
    steps JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB,
    UNIQUE (name, version)
);

CREATE INDEX IF NOT EXISTS idx_workflow_definitions_is_enabled ON workflow_definitions (is_enabled);
CREATE INDEX IF NOT EXISTS idx_workflow_definitions_created_at ON workflow_definitions (created_at);
CREATE INDEX IF NOT EXISTS idx_workflow_definitions_updated_at ON workflow_definitions (updated_at);