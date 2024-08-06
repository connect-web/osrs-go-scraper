CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;
CREATE EXTENSION pg_prewarm;


-- Prewarm the table into the buffer cache
-- SELECT pg_prewarm('your_table_name');