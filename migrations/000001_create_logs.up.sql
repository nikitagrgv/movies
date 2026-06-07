CREATE SCHEMA IF NOT EXISTS logs;

CREATE TABLE IF NOT EXISTS logs.visits
(
    id           INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    ip_address   INET        NOT NULL,
    path         TEXT        NOT NULL,
    duration     BIGINT      NOT NULL,
    attempted_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_visits_attempted_at ON logs.visits (attempted_at DESC);
CREATE INDEX IF NOT EXISTS idx_visits_ip_address ON logs.visits (ip_address);
