CREATE EXTENSION pg_prewarm;
SELECT pg_prewarm('not_found');
SELECT pg_prewarm('ml.results');
SELECT pg_prewarm('grouped.skillers');
SELECT pg_prewarm('ml.metrics');
SELECT pg_prewarm('ml.metrics_large');
SELECT pg_prewarm('users.accs');
SELECT pg_prewarm('users.fiber_storage');