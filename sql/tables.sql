-- The players table storing player_id's, Name's and first appearances to the database
CREATE TABLE IF NOT EXISTS players(
    ID SERIAL PRIMARY KEY,
    NAME VARCHAR(12) UNIQUE NOT NULL,
    FIRST_SEEN timestamptz NOT NULL DEFAULT NOW()
);


-- The user is not found on Hiscores
-- This can mean: Banned || Username was changed
CREATE TABLE IF NOT EXISTS not_found(
   PlayerId SERIAL PRIMARY KEY,
   FIRST_SEEN timestamptz NOT NULL DEFAULT NOW(),
   FOREIGN KEY (PlayerId) REFERENCES players(id)
);

-- Low Latency Live player stats
CREATE TABLE IF NOT EXISTS player_live(
    PlayerId INT PRIMARY KEY,
    LAST_UPDATED timestamptz not null default NOW(),
    skills_experience JSONB,
    skills_levels JSONB,
    skills_ratio JSONB,
    minigames JSONB,
    FOREIGN KEY (PlayerId) REFERENCES players(id)
);