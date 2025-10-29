CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS workflow_definitions
(
    id          UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    description TEXT,
    version     TEXT NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_workflow_name_version
    ON workflow_definitions (name, version);