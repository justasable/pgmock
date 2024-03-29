BEGIN;

-- public schema
CREATE TABLE public.Types (
    type_int integer,
    type_bool bool,
    type_numeric numeric,
    type_numericp numeric(10, 2),
    type_text text,
    type_timestamptz timestamptz,
    type_date date,
    type_byte bytea,
    type_uuid uuid
);

CREATE TABLE public.Identity (
    identity_always integer GENERATED ALWAYS AS IDENTITY,
    identity_default integer GENERATED BY DEFAULT AS IDENTITY,
    PRIMARY KEY (identity_always)
);

CREATE TABLE public.Constraints (
    con_pk_one integer,
    con_pk_two integer,
    con_null integer,
    con_not_null integer NOT NULL,
    con_default integer DEFAULT 4,
    con_no_default integer,
    con_generated integer GENERATED ALWAYS AS (7) STORED,
    PRIMARY KEY (con_pk_one, con_pk_two)
);

-- public view
CREATE VIEW public.hello AS
SELECT 'hello world';

-- test schema
CREATE SCHEMA test;

CREATE TABLE test.References (
    fk_single integer,
    fk_multiple_one integer,
    fk_multiple_two integer,
    FOREIGN KEY (fk_single) REFERENCES public.Identity (identity_always),
    FOREIGN KEY (fk_multiple_one, fk_multiple_two) REFERENCES public.Constraints(con_pk_one, con_pk_two)
);

COMMIT;