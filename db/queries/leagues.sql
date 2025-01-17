-- name: GetLeagueById :one
SELECT "id", "leagueId", "year", "teamCount", "currentWeek", "nflWeek" 
FROM "leagues" 
WHERE "id" = $1;

-- name: GetLeagues :many
SELECT "id", "leagueId", "year", "teamCount", "currentWeek", "nflWeek"
FROM "leagues"
WHERE ("year" = $1)
AND ("teamCount" = $2 OR $2 IS NULL)
AND ("currentWeek" = $3 OR $3 IS NULL)
AND ("nflWeek" = $4 OR $4 IS NULL)
ORDER BY
CASE $5
    WHEN 'id' THEN "id"
    WHEN 'year' THEN "year"
    WHEN 'teamCount' THEN "teamCount"
END,
CASE 
    WHEN $5 LIKE '-%' THEN 'DESC'
    ELSE 'ASC'
END;


