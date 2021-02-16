CREATE TABLE blocks
(
    id                      BIGSERIAL,
    block_hash              CHAR(64) NOT NULL,
    blue_score              BIGINT CHECK (blue_score >= 0) NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT idx_blocks_block_hash UNIQUE  (block_hash)
);