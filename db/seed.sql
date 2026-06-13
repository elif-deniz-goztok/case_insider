-- Teams (strength values reflect real-world relative quality)
INSERT INTO teams (name, strength) VALUES
    ('Chelsea',          9),
    ('Manchester City',  8),
    ('Arsenal',          7),
    ('Liverpool',        6)
ON CONFLICT (name) DO NOTHING;

-- Fixtures: double round-robin, 6 weeks x 2 matches
-- Week 1
INSERT INTO matches (week, home_team_id, away_team_id)
SELECT 1, t1.id, t2.id FROM teams t1, teams t2 WHERE t1.name = 'Chelsea'          AND t2.name = 'Arsenal';
INSERT INTO matches (week, home_team_id, away_team_id)
SELECT 1, t1.id, t2.id FROM teams t1, teams t2 WHERE t1.name = 'Manchester City'  AND t2.name = 'Liverpool';

-- Week 2
INSERT INTO matches (week, home_team_id, away_team_id)
SELECT 2, t1.id, t2.id FROM teams t1, teams t2 WHERE t1.name = 'Chelsea'          AND t2.name = 'Manchester City';
INSERT INTO matches (week, home_team_id, away_team_id)
SELECT 2, t1.id, t2.id FROM teams t1, teams t2 WHERE t1.name = 'Arsenal'          AND t2.name = 'Liverpool';

-- Week 3
INSERT INTO matches (week, home_team_id, away_team_id)
SELECT 3, t1.id, t2.id FROM teams t1, teams t2 WHERE t1.name = 'Chelsea'          AND t2.name = 'Liverpool';
INSERT INTO matches (week, home_team_id, away_team_id)
SELECT 3, t1.id, t2.id FROM teams t1, teams t2 WHERE t1.name = 'Arsenal'          AND t2.name = 'Manchester City';

-- Week 4 (reversed home/away from week 1)
INSERT INTO matches (week, home_team_id, away_team_id)
SELECT 4, t1.id, t2.id FROM teams t1, teams t2 WHERE t1.name = 'Arsenal'          AND t2.name = 'Chelsea';
INSERT INTO matches (week, home_team_id, away_team_id)
SELECT 4, t1.id, t2.id FROM teams t1, teams t2 WHERE t1.name = 'Liverpool'        AND t2.name = 'Manchester City';

-- Week 5 (reversed home/away from week 2)
INSERT INTO matches (week, home_team_id, away_team_id)
SELECT 5, t1.id, t2.id FROM teams t1, teams t2 WHERE t1.name = 'Manchester City'  AND t2.name = 'Chelsea';
INSERT INTO matches (week, home_team_id, away_team_id)
SELECT 5, t1.id, t2.id FROM teams t1, teams t2 WHERE t1.name = 'Liverpool'        AND t2.name = 'Arsenal';

-- Week 6 (reversed home/away from week 3)
INSERT INTO matches (week, home_team_id, away_team_id)
SELECT 6, t1.id, t2.id FROM teams t1, teams t2 WHERE t1.name = 'Liverpool'        AND t2.name = 'Chelsea';
INSERT INTO matches (week, home_team_id, away_team_id)
SELECT 6, t1.id, t2.id FROM teams t1, teams t2 WHERE t1.name = 'Manchester City'  AND t2.name = 'Arsenal';
