CREATE TABLE universities (
    slug TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    color_hex TEXT NOT NULL DEFAULT '#888888'
);

CREATE TABLE faculties (
    slug TEXT PRIMARY KEY,
    university_slug TEXT NOT NULL REFERENCES universities(slug),
    name TEXT NOT NULL,
    short_abbr TEXT NOT NULL,
    color_hex TEXT NOT NULL DEFAULT '#888888'
);

CREATE TABLE courses (
    slug TEXT PRIMARY KEY,
    faculty_slug TEXT NOT NULL REFERENCES faculties(slug),
    name TEXT NOT NULL,
    short_abbr TEXT NOT NULL,
    color_hex TEXT NOT NULL DEFAULT '#888888'
);

ALTER TABLE players
    ADD COLUMN affiliation_type TEXT NOT NULL DEFAULT 'student',
    ADD COLUMN university_slug TEXT,
    ADD COLUMN faculty_slug TEXT,
    ADD COLUMN course_slug TEXT;

ALTER TABLE invites
    ADD COLUMN affiliation_type TEXT NOT NULL DEFAULT 'student';

CREATE INDEX idx_players_affiliation_type ON players (affiliation_type);

INSERT INTO universities (slug, name, color_hex) VALUES
    ('campus-demo', 'CampusWorld Demo', '#5eead4');

INSERT INTO faculties (slug, university_slug, name, short_abbr, color_hex) VALUES
    ('ccet', 'campus-demo', 'Centro de Ciências Exatas e Tecnologia', 'CCET', '#60a5fa'),
    ('chs', 'campus-demo', 'Centro de Ciências Humanas e Sociais', 'CHS', '#f472b6'),
    ('health', 'campus-demo', 'Centro de Ciências da Saúde', 'SAU', '#34d399');

INSERT INTO courses (slug, faculty_slug, name, short_abbr, color_hex) VALUES
    ('ads', 'ccet', 'Análise e Desenvolvimento de Sistemas', 'ADS', '#38bdf8'),
    ('eng-comp', 'ccet', 'Engenharia da Computação', 'EC', '#2563eb'),
    ('pedagogia', 'chs', 'Pedagogia', 'PED', '#ec4899'),
    ('medicina', 'health', 'Medicina', 'MED', '#10b981');
