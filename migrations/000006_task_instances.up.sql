CREATE TABLE IF NOT EXISTS step_instances
(
    id                   UUID PRIMARY KEY,
    workflow_instance_id UUID        NOT NULL,
    step_id              VARCHAR     NOT NULL,
    status               VARCHAR(50) NOT NULL,
    started_at           TIMESTAMPTZ,
    completed_at         TIMESTAMPTZ,
    input                JSONB,
    output               JSONB,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata             JSONB,
    CONSTRAINT fk_workflow_instance
        FOREIGN KEY (workflow_instance_id)
            REFERENCES workflow_instances (id)
            ON DELETE CASCADE
);