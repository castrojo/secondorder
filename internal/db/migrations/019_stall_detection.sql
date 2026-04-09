-- SO-87: stall detection threshold setting
-- Ensure settings table exists (guard for legacy test schemas that skip earlier migrations).
CREATE TABLE IF NOT EXISTS settings (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
-- Default: 4 hours (= 2 sessions) per AC3
INSERT OR IGNORE INTO settings (key, value) VALUES ('stall_threshold_hours', '4');
