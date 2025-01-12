-- name: GetLeagueById :one
SELECT "id", "leagueId", "year", "teamCount", "currentWeek", "nflWeek" 
FROM "leagues" 
WHERE "id" = $1;

