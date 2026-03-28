-- Work block indexes for lifecycle queries
CREATE INDEX IF NOT EXISTS idx_work_blocks_status ON work_blocks(status);
CREATE INDEX IF NOT EXISTS idx_issues_work_block ON issues(work_block_id);
