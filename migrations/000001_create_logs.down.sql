DROP INDEX IF EXISTS logs.idx_visits_attempted_at;
DROP INDEX IF EXISTS logs.idx_visits_ip_address;

DROP TABLE IF EXISTS logs.visits;

DROP SCHEMA IF EXISTS logs;