-- CREATE TABLE "leagues" (
--     "id" INTEGER PRIMARY KEY,
--     "leagueId" INTEGER NOT NULL,
--     "year" INTEGER NOT NULL,
--     "teamCount" INTEGER,
--     "currentWeek" INTEGER,
--     "nflWeek" INTEGER,
--     CONSTRAINT "uix_league_year" UNIQUE ("leagueId", "year")
-- );

CREATE TABLE leagues (
    "id" SERIAL PRIMARY KEY,
    "leagueId" INTEGER NOT NULL,
    "year" INTEGER NOT NULL,
    "teamCount" INTEGER NOT NULL,
    "currentWeek" INTEGER NOT NULL DEFAULT 0,
    "nflWeek" INTEGER NOT NULL DEFAULT 0,
    
    CONSTRAINT "uix_league_year" UNIQUE ("leagueId", "year")
);

-- CREATE TABLE settings (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     league_id INTEGER NOT NULL UNIQUE,
--     regularSeasonCount INTEGER,
--     vetoVotesRequired INTEGER,
--     teamCount INTEGER,
--     playoffTeamCount INTEGER,
--     keeperCount INTEGER,
--     tradeDeadline BIGINT,
--     name VARCHAR(255),
--     tieRule VARCHAR(50),
--     playoffTieRule VARCHAR(50),
--     playoffSeedTieRule VARCHAR(50),
--     playoffMatchupPeriodLength INTEGER,
--     faab BOOLEAN,
--     CONSTRAINT idx_settings_league INDEX (league_id),
--     FOREIGN KEY (league_id) REFERENCES leagues(id)
-- );

CREATE TABLE "teams" (
    "id" INTEGER PRIMARY KEY,
    "league_id" INTEGER NOT NULL,
    "teamId" INTEGER NOT NULL,
    "year" INTEGER NOT NULL,
    "teamAbbrv" VARCHAR(10) NOT NULL,
    "teamName" VARCHAR(255) NOT NULL,
    "owners" VARCHAR(50),
    "divisionId" VARCHAR(255),
    "divisionName" VARCHAR(255),
    "wins" INTEGER DEFAULT 0,
    "losses" INTEGER DEFAULT 0,
    "ties" INTEGER DEFAULT 0,
    "pointsFor" INTEGER DEFAULT 0,
    "pointsAgainst" INTEGER DEFAULT 0,
    "waiverRank" INTEGER,
    "acquisitions" INTEGER DEFAULT 0,
    "acquisitionBudgetSpent" INTEGER DEFAULT 0,
    "drops" INTEGER DEFAULT 0,
    "trades" INTEGER DEFAULT 0,
    "streakType" VARCHAR(50),
    "streakLength" INTEGER,
    "standing" INTEGER,
    "finalStanding" INTEGER,
    "draftProjRank" INTEGER,
    "playoffPct" INTEGER,
    "logoUrl" VARCHAR(255),
    CONSTRAINT "uix_team_year" UNIQUE ("teamId", "year"),
    FOREIGN KEY ("league_id") REFERENCES "leagues"("id")
);


-- CREATE TABLE players (
--     id INTEGER PRIMARY KEY,
--     espnId INTEGER UNIQUE NOT NULL,
--     name VARCHAR(255) NOT NULL,
--     position VARCHAR(50),
--     CONSTRAINT idx_player_name INDEX (name),
--     CONSTRAINT idx_player_position INDEX (position)
-- );

-- CREATE TABLE drafts (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     team_id INTEGER NOT NULL,
--     player_id INTEGER NOT NULL,
--     overallPick INTEGER NOT NULL,
--     roundNum INTEGER NOT NULL,
--     roundPick INTEGER NOT NULL,
--     keeperStatus BOOLEAN DEFAULT FALSE,
--     bidAmount INTEGER DEFAULT NULL,
--     nominating_team_id INTEGER DEFAULT NULL,
--     CONSTRAINT uix_draft_pick UNIQUE (team_id, player_id),
--     CONSTRAINT idx_draft_team_player_year INDEX (team_id, player_id),
--     FOREIGN KEY (team_id) REFERENCES teams(id),
--     FOREIGN KEY (player_id) REFERENCES players(id),
--     FOREIGN KEY (nominating_team_id) REFERENCES teams(id)
-- );

-- CREATE TABLE matchups (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     week INTEGER NOT NULL,
--     home_team_id INTEGER,
--     away_team_id INTEGER,
--     homeScore FLOAT,
--     awayScore FLOAT,
--     isPlayoff BOOLEAN DEFAULT FALSE,
--     matchupType VARCHAR(50) DEFAULT 'NONE',
--     CONSTRAINT uix_matchup UNIQUE (week, home_team_id, away_team_id),
--     CONSTRAINT idx_matchup_team_week INDEX (home_team_id, away_team_id, week),
--     CONSTRAINT idx_matchup_week INDEX (week),
--     FOREIGN KEY (home_team_id) REFERENCES teams(id),
--     FOREIGN KEY (away_team_id) REFERENCES teams(id)
-- );

-- CREATE TABLE activities (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     date BIGINT,
--     team_id INTEGER,
--     player_id INTEGER NOT NULL,
--     bidAmount FLOAT,
--     action VARCHAR(50),
--     CONSTRAINT idx_activity_team_player INDEX (team_id, player_id),
--     CONSTRAINT idx_activity_team INDEX (team_id),
--     FOREIGN KEY (team_id) REFERENCES teams(id),
--     FOREIGN KEY (player_id) REFERENCES players(id)
-- );

-- CREATE TABLE rosters (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     team_id INTEGER NOT NULL,
--     player_id INTEGER NOT NULL,
--     rosterSlot VARCHAR(50),
--     CONSTRAINT uix_roster_team_player UNIQUE (team_id, player_id),
--     CONSTRAINT idx_roster_team_player INDEX (team_id, player_id),
--     CONSTRAINT idx_roster_team INDEX (team_id),
--     FOREIGN KEY (team_id) REFERENCES teams(id),
--     FOREIGN KEY (player_id) REFERENCES players(id)
-- );

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    email citext,
    activated bool NOT NULL,
    provider text NOT NULL,
    provider_id text NOT NULL,
    version integer NOT NULL DEFAULT 1,
    CONSTRAINT unique_provider_and_provider_id UNIQUE (provider, provider_id)
);

CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL
);