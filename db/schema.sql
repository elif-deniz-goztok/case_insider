CREATE TABLE IF NOT EXISTS teams (
    id       SERIAL PRIMARY KEY,
    name     VARCHAR(100) NOT NULL UNIQUE,
    strength INTEGER NOT NULL CHECK (strength BETWEEN 1 AND 10)
);

CREATE TABLE IF NOT EXISTS matches (
    id           SERIAL PRIMARY KEY,
    week         INTEGER NOT NULL,
    home_team_id INTEGER NOT NULL REFERENCES teams(id),
    away_team_id INTEGER NOT NULL REFERENCES teams(id),
    home_goals   INTEGER,
    away_goals   INTEGER,
    played       BOOLEAN NOT NULL DEFAULT FALSE
);
