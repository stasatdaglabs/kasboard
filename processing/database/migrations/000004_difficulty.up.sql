ALTER TABLE blocks ADD COLUMN difficulty DOUBLE PRECISION;
UPDATE blocks SET difficulty = 0;
ALTER TABLE blocks ALTER COLUMN difficulty SET NOT NULL;

CREATE TABLE estimated_blue_hashrates
(
    timestamp     BIGINT NOT NULL,
    blue_hashrate BIGINT NOT NULL
);
