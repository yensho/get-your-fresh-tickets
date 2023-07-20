CREATE SCHEMA gyft;

CREATE TABLE spaces (
    space_nm text NOT NULL PRIMARY KEY,
    space_section_nm text,
    space_section_seats blob
);

CREATE TABLE auth (
    user text NOT NULL PRIMARY KEY,
    token text NOT NULL
)`