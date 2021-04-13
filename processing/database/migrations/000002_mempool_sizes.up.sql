CREATE TABLE mempool_sizes
(
    timestamp BIGINT NOT NULL,
    size      BIGINT NOT NULL
);

CREATE TABLE pruning_point_movements
(
    timestamp                BIGINT   NOT NULL,
    pruning_point_block_hash CHAR(64) NOT NULL
);
