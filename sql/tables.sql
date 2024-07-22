-- The players table storing player_id's, Name's and first appearances to the database
CREATE TABLE IF NOT EXISTS players
(
    id         SERIAL PRIMARY KEY,
    name       CITEXT UNIQUE NOT NULL,
    first_seen timestamptz   NOT NULL DEFAULT NOW()
);

-- These players do not need to be scraped they are low priority.
CREATE TABLE IF NOT EXISTS old_players
(
    id         SERIAL,
    name       CITEXT PRIMARY KEY NOT NULL,
    first_seen timestamptz        NOT NULL DEFAULT NOW()
);

-- The user is not found on Hiscores
-- This can mean: Banned || Username was changed
CREATE TABLE IF NOT EXISTS not_found
(
    PlayerId   SERIAL PRIMARY KEY,
    first_seen timestamptz NOT NULL DEFAULT NOW(),
    FOREIGN KEY (PlayerId) REFERENCES players (id)
);

-- Low Latency Live player stats
CREATE TABLE IF NOT EXISTS player_live
(
    PlayerId          INT PRIMARY KEY,
    last_updated      timestamptz not null default NOW(),
    skills_experience JSONB,
    skills_levels     JSONB,
    skills_ratio      JSONB,
    minigames         JSONB,
    FOREIGN KEY (PlayerId) REFERENCES players (id)
);


CREATE TABLE IF NOT EXISTS player_live_stats
(
    PlayerId     INT PRIMARY KEY,
    LAST_UPDATED timestamptz not null, -- THIS MUST MATCH THE player_live LAST_UPDATED
    combat_level smallint,
    Overall      bigint,
    total_level  smallint,
    FOREIGN KEY (PlayerId) REFERENCES players (id)
);

CREATE TABLE IF NOT EXISTS player_gains
(
    PlayerId     INT PRIMARY KEY,
    LAST_UPDATED timestamptz not null DEFAULT NOW(), -- THIS MUST MATCH THE player_live LAST_UPDATED
    skills_experience JSONB,
    skills_ratio      JSONB,
    minigames         JSONB,
    FOREIGN KEY (PlayerId) REFERENCES players (id)
);


CREATE SCHEMA IF NOT EXISTS STATS;

CREATE TABLE IF NOT EXISTS stats.Pearson(
    PlayerId     INT PRIMARY KEY,
    LAST_UPDATED timestamptz not null, -- THIS MUST MATCH THE player_live LAST_UPDATED
    skill TEXT,
    linked_players INTEGER[],
    FOREIGN KEY (PlayerId) REFERENCES players (id)
);