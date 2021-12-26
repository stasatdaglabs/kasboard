ALTER TABLE blocks ADD COLUMN difficulty DOUBLE PRECISION;
UPDATE blocks SET difficulty = 0;
ALTER TABLE blocks ALTER COLUMN difficulty SET NOT NULL;
