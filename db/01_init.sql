DROP DATABASE postgres;

\connect

CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;

CREATE TABLE IF NOT EXISTS public.music_groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);
CREATE INDEX ON public.music_groups USING btree (name);

CREATE TABLE IF NOT EXISTS public.songs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    group_id INT REFERENCES public.music_groups(id) NOT NULL,
    release_date VARCHAR(10) NOT NULL,
    text TEXT[] NOT NULL,
    link VARCHAR(255) NOT NULL,
    deleted TIMESTAMP
);
CREATE INDEX ON public.songs USING btree (name);
