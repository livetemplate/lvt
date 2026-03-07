CREATE TABLE IF NOT EXISTS generation_events (
    id TEXT PRIMARY KEY,
    timestamp DATETIME NOT NULL,
    command TEXT NOT NULL,
    inputs TEXT NOT NULL,
    kit TEXT,
    lvt_version TEXT,
    success BOOLEAN NOT NULL,
    validation TEXT,
    errors TEXT,
    duration_ms INTEGER,
    files_generated TEXT,
    components_used TEXT,
    component_errors TEXT
);
CREATE INDEX IF NOT EXISTS idx_events_timestamp ON generation_events(timestamp);
CREATE INDEX IF NOT EXISTS idx_events_success ON generation_events(success);
CREATE INDEX IF NOT EXISTS idx_events_command ON generation_events(command);
CREATE INDEX IF NOT EXISTS idx_events_kit ON generation_events(kit);
