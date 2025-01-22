-- name: GetLeagueById :one
SELECT "id", "leagueId", "year", "teamCount", "currentWeek", "nflWeek" 
FROM "leagues" 
WHERE "id" = $1;

-- name: GetLeagues :many
SELECT
    "id",
    "leagueId",
    "year",
    "teamCount",
    "currentWeek",
    "nflWeek"
FROM
    "leagues"
WHERE
    ("id" = COALESCE($1, "id") OR $1 IS NULL)
    AND ("leagueId" = COALESCE($2, "leagueId") OR $2 IS NULL)
    AND ("year" = COALESCE($3, "year") OR $3 IS NULL)
    AND ("teamCount" = COALESCE($4, "teamCount") OR $4 IS NULL)
    AND ("currentWeek" = COALESCE($5, "currentWeek") OR $5 IS NULL)
    AND ("nflWeek" = COALESCE($6, "nflWeek") OR $6 IS NULL)
ORDER BY
    CASE 
        WHEN $7 = 'id' THEN "id"
        WHEN $7 = 'year' THEN "year"
        WHEN $7 = 'teamCount' THEN "teamCount"
        WHEN $7 = 'currentWeek' THEN "currentWeek"
        WHEN $7 = 'nflWeek' THEN "nflWeek"
        ELSE "id"
    END,
    CASE
        WHEN $7 LIKE '-%' THEN 'DESC'
        ELSE 'ASC'
    END
LIMIT $8
OFFSET $9;







