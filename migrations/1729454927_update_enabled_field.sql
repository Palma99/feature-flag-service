ALTER TABLE flag_environment ADD COLUMN enabled boolean NOT NULL DEFAULT false;
ALTER TABLE flag DROP COLUMN enabled;