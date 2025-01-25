-- name: GetTeamById :one
SELECT "id", "league_id", "year", "teamAbbrv", "owners", "divisionId", "divisionName", "wins", "losses", "ties", "pointsFor", "pointsAgainst", "waiverRank", "acquisitions", "acquisitionBudgetSpent", "drops", "trades", "logoUrl" 
FROM "teams" 
WHERE "id" = $1;