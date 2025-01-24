-- name: GetLeagueById :one
SELECT "id", "leagueId", "year", "teamCount", "currentWeek", "nflWeek" 
FROM "leagues" 
WHERE "id" = $1;

-- -- name: GetLeagues :many
-- SELECT
--     "id",
--     "leagueId",
--     "year",
--     "teamCount",
--     "currentWeek",
--     "nflWeek"
-- FROM
--     leagues
-- WHERE
--     ("id" = $1 OR $1 = -1)
--     AND ("leagueId" = $2 OR $2 = -1)
--     AND ("year" = $3 OR $3 = -1)
--     AND ("teamCount" = $4 OR $4 = -1)
--     AND ("currentWeek" = $5 OR $5 = -1)
--     AND ("nflWeek" = $6 OR $6 = -1)
-- LIMIT $7
-- OFFSET $8;

-- name: GetLeaguesAsc :many
SELECT
    "id",
    "leagueId",
    "year",
    "teamCount",
    "currentWeek",
    "nflWeek"
FROM
    leagues
WHERE
    ("id" = $1 OR $1 = -1)
    AND ("leagueId" = $2 OR $2 = -1)
    AND ("year" = $3 OR $3 = -1)
    AND ("teamCount" = $4 OR $4 = -1)
    AND ("currentWeek" = $5 OR $5 = -1)
    AND ("nflWeek" = $6 OR $6 = -1)
ORDER BY
    CASE
        WHEN $9 = 'id' THEN "id"
        WHEN $9 = 'leagueId' THEN "leagueId"
        WHEN $9 = 'year' THEN "year"
        WHEN $9 = 'teamCount' THEN "teamCount"
        WHEN $9 = 'currentWeek' THEN "currentWeek"
        WHEN $9 = 'nflWeek' THEN "nflWeek"
        ELSE "id"
    END ASC
LIMIT $7
OFFSET $8;

-- name: GetLeaguesDesc :many
SELECT
    "id",
    "leagueId",
    "year",
    "teamCount",
    "currentWeek",
    "nflWeek"
FROM
    leagues
WHERE
    ("id" = $1 OR $1 = -1)
    AND ("leagueId" = $2 OR $2 = -1)
    AND ("year" = $3 OR $3 = -1)
    AND ("teamCount" = $4 OR $4 = -1)
    AND ("currentWeek" = $5 OR $5 = -1)
    AND ("nflWeek" = $6 OR $6 = -1)
ORDER BY
    CASE
        WHEN $9 = 'id' THEN "id"
        WHEN $9 = 'leagueId' THEN "leagueId"
        WHEN $9 = 'year' THEN "year"
        WHEN $9 = 'teamCount' THEN "teamCount"
        WHEN $9 = 'currentWeek' THEN "currentWeek"
        WHEN $9 = 'nflWeek' THEN "nflWeek"
        ELSE "id"
    END DESC
LIMIT $7
OFFSET $8;









