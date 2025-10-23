CREATE TABLE IF NOT EXISTS workflow_instances
(
    id                     UUID PRIMARY KEY,
    workflow_definition_id UUID        NOT NULL,
    status                 VARCHAR(50) NOT NULL,
    started_at             TIMESTAMPTZ,
    completed_at           TIMESTAMPTZ,
    input                  JSONB,
    output                 JSONB,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata               JSONB,
    CONSTRAINT fk_workflow_definition
        FOREIGN KEY (workflow_definition_id)
            REFERENCES workflow_definitions (id)
            ON DELETE CASCADE
);