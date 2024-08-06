-- The players table storing player_id's, Name's and first appearances to the database
CREATE TABLE IF NOT EXISTS players
(
    id         SERIAL PRIMARY KEY,
    name       CITEXT UNIQUE NOT NULL,
    first_seen timestamptz   NOT NULL DEFAULT NOW()
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


CREATE TABLE IF NOT EXISTS player_live_stats(
    PlayerId     INT PRIMARY KEY,
    LAST_UPDATED timestamptz not null, -- THIS MUST MATCH THE player_live LAST_UPDATED
    combat_level smallint DEFAULT 3,
    Overall      bigint DEFAULT 0,
    total_level  smallint DEFAULT 32,
    FOREIGN KEY (PlayerId) REFERENCES players (id)
);

CREATE TABLE IF NOT EXISTS player_gains
(
    PlayerId          INT PRIMARY KEY,
    LAST_UPDATED      timestamptz not null DEFAULT NOW(), -- THIS MUST MATCH THE player_live LAST_UPDATED
    skills_experience JSONB,
    skills_ratio      JSONB,
    minigames         JSONB,
    FOREIGN KEY (PlayerId) REFERENCES players (id)
);



CREATE SCHEMA IF NOT EXISTS ML;

CREATE TABLE IF NOT EXISTS ML.results
(
    PlayerId     int,
    Activity     TEXT,
    ActivityType TEXT,
    TIME         timestamptz,
    PRIMARY KEY (PlayerId, Activity),
    FOREIGN KEY (PlayerId) REFERENCES players (id)
);

CREATE TABLE IF NOT EXISTS ML.results_large
(
    PlayerId     int,
    Activity     TEXT,
    ActivityType TEXT,
    TIME         timestamptz,
    PRIMARY KEY (PlayerId, Activity),
    FOREIGN KEY (PlayerId) REFERENCES players (id)
);



CREATE TABLE IF NOT EXISTS stats.bulk_minigames
(
    ID       INT PRIMARY KEY,
    minigame TEXT,
    links    INTEGER[],
    FOREIGN KEY (ID) REFERENCES PLAYERS (ID)
);



CREATE SCHEMA IF NOT EXISTS GROUPED;

CREATE TABLE IF NOT EXISTS grouped.skills
(
    ID           SERIAL PRIMARY KEY,
    skills       TEXT[] UNIQUE,
    skills_index INTEGER[] UNIQUE,
    skills_csv   TEXT UNIQUE
);



CREATE TABLE IF NOT EXISTS grouped.skills_count
(
    ID     INT PRIMARY KEY,
    Amount int,
    FOREIGN KEY (id) references grouped.skills (id)
);

CREATE TABLE IF NOT EXISTS grouped.skills_users
(
    ID        INT PRIMARY KEY,
    PlayerIds INTEGER[],
    FOREIGN KEY (id) references grouped.skills (id)
);



CREATE TABLE IF NOT EXISTS stats.minigame_links
(
    ID       INT PRIMARY KEY,
    minigame TEXT,
    links    INTEGER,
    FOREIGN KEY (ID) REFERENCES PLAYERS (ID)
);

CREATE TABLE IF NOT EXISTS stats.minigame_links
(
    ID       INT PRIMARY KEY,
    minigame TEXT,
    links    INTEGER[],
    FOREIGN KEY (ID) REFERENCES PLAYERS (ID)
);



CREATE TABLE IF NOT EXISTS ML.results_large
(
    PlayerId     int,
    Activity     TEXT,
    ActivityType TEXT,
    TIME         timestamptz,
    PRIMARY KEY (PlayerId, Activity),
    FOREIGN KEY (PlayerId) REFERENCES players (id)
);


CREATE TABLE IF NOT EXISTS ML.metrics_large
(
    Activity TEXT PRIMARY KEY,
    metrics  JSONB,
    TIME     timestamptz
);


CREATE SCHEMA IF NOT EXISTS GROUPED;

CREATE TABLE IF NOT EXISTS grouped.skillers
(
    id        INTEGER PRIMARY KEY,
    skills    TEXT[],
    amount    int,
    playerIds INTEGER[]
);