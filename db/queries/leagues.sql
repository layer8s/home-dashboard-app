-- name: GetLeagueById :one
SELECT "id", "leagueId", "year", "teamCount", "currentWeek", "nflWeek" 
FROM "leagues" 
WHERE "id" = $1;

-- name: GetLeagues :many
SELECT "id", "leagueId", "year", "teamCount", "currentWeek", "nflWeek" 
FROM "leagues";


