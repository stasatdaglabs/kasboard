ALTER TABLE blocks ADD COLUMN difficulty DOUBLE PRECISION;
UPDATE blocks SET difficulty = 0;
ALTER TABLE blocks ALTER COLUMN difficulty SET NOT NULL;

ALTER TABLE blocks ADD COLUMN propagation_delay DOUBLE PRECISION;
UPDATE blocks SET propagation_delay = 0;
ALTER TABLE blocks ALTER COLUMN propagation_delay SET NOT NULL;

ALTER TABLE analyzed_blocks ADD COLUMN average_propagation_delay DOUBLE PRECISION;
UPDATE analyzed_blocks SET average_propagation_delay = 0;
ALTER TABLE analyzed_blocks ALTER COLUMN average_propagation_delay SET NOT NULL;

CREATE TABLE estimated_blue_hashrates
(
    timestamp     BIGINT NOT NULL,
    blue_hashrate BIGINT NOT NULL
);
