CREATE TABLE blocks
(
    id         BIGSERIAL,
    block_hash CHAR(64)                                        NOT NULL,
    timestamp  BIGINT                                          NOT NULL,
    blue_score BIGINT CHECK (blue_score >= 0)                  NOT NULL,
    bits       BIGINT CHECK (bits >= 0 AND bits <= 4294967295) NOT NULL,
    hashrate   BIGINT                                          NOT NULL,
    PRIMARY KEY (id)
);