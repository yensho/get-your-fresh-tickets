CREATE SCHEMA gyft;

CREATE TABLE gyft.customers (
    id        uuid NOT NULL PRIMARY KEY,
    cust_name text NOT NULL,
    age       int,
    gender    text,
    phone_num text,
    addr      jsonb,
    email     text,
    isrt_ts   timestamptz,
    isrt_usr  text,
    updt_ts   timestamptz,
    updt_usr  text
);

CREATE FUNCTION cust_ins() RETURNS trigger AS $cust_ins$
    BEGIN
        NEW.id := gen_random_uuid();
        NEW.isrt_ts := current_timestamp;
        NEW.isrt_usr := current_user;
        RETURN NEW;
    END;
$cust_ins$ LANGUAGE plpgsql;

CREATE TRIGGER cust_ins BEFORE INSERT ON gyft.customers
    FOR EACH ROW EXECUTE FUNCTION cust_ins();

CREATE TABLE gyft.shows (
    id           uuid NOT NULL PRIMARY KEY,
    show_name    text,
    loc          text,
    show_times   timestamptz,
    seats        text[],
    prices       jsonb,
    age_restrict boolean,
    isrt_ts      timestamptz,
    isrt_usr     text,
    updt_ts      timestamptz,
    upst_usr     text
);

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