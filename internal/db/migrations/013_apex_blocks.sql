-- Migration 013: Apex Blocks and North Star metrics
CREATE TABLE IF NOT EXISTS apex_blocks (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    goal TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

ALTER TABLE work_blocks ADD COLUMN north_star_metric TEXT NOT NULL DEFAULT '';
ALTER TABLE work_blocks ADD COLUMN north_star_target TEXT NOT NULL DEFAULT '';
ALTER TABLE work_blocks ADD COLUMN apex_block_id TEXT REFERENCES apex_blocks(id);
