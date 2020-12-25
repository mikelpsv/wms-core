--
-- PostgreSQL database dump
--

-- Dumped from database version 13.1
-- Dumped by pg_dump version 13.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: storage1; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.storage1 (
    zone_id integer,
    cell_id integer,
    prod_id integer,
    quantity integer
);


ALTER TABLE public.storage1 OWNER TO devuser;

--
-- Data for Name: storage1; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.storage1 (zone_id, cell_id, prod_id, quantity) FROM stdin;
0	2	32	100
0	2	34	40
0	2	34	40
0	2	34	40
0	2	34	40
0	2	34	40
0	2	32	100
\.


--
-- PostgreSQL database dump complete
--

