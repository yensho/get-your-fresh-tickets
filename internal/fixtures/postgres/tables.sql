CREATE SCHEMA gyft;

CREATE TABLE gyft.spaces (
    space_nm text NOT NULL PRIMARY KEY,
    space_section_nm text,
    space_section_seats jsonb
);

CREATE TABLE gyft.auth (
    user_name text NOT NULL PRIMARY KEY,
    pass_token text NOT NULL
);

CREATE USER gyftusr WITH PASSWORD 'gyftPwE0';
GRANT ALL ON SCHEMA gyft TO gyftusr;
GRANT ALL ON ALL TABLES IN SCHEMA gyft TO gyftusr;