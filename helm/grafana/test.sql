WITH fresh AS (
	SELECT *
	FROM logbeast.db.freshmeat
	WHERE timestamp >= CURRENT_DATE - INTERVAL '7' DAY
),
cured AS (
	SELECT *
	FROM logbeast.db.hive
	WHERE timestamp < CURRENT_DATE - INTERVAL '7' DAY
),
combined AS (
	SELECT *
	FROM fresh

	UNION ALL

	SELECT *
	FROM cured
),
hive AS (
	SELECT *
	FROM combined
	WHERE timestamp >= FROM_ISO8601_TIMESTAMP('2026-07-08T22:00:00+00:00')
		AND timestamp < FROM_ISO8601_TIMESTAMP('2026-07-09T22:38:25+00:00')
)
SELECT count(*) FROM hive;
