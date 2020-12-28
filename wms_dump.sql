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
-- Name: cells; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.cells (
    id integer NOT NULL,
    name character varying(50) DEFAULT ''::character varying NOT NULL,
    whs_id integer DEFAULT 0 NOT NULL,
    zone_id integer DEFAULT 0 NOT NULL,
    passage_id integer DEFAULT 0 NOT NULL,
    rack_id integer DEFAULT 0 NOT NULL,
    floor integer DEFAULT 0 NOT NULL,
    size_length integer DEFAULT 0 NOT NULL,
    size_width integer DEFAULT 0 NOT NULL,
    size_height integer DEFAULT 0 NOT NULL,
    size_volume integer DEFAULT 0 NOT NULL,
    size_usefull_volume integer DEFAULT 0 NOT NULL,
    size_weight integer DEFAULT 0 NOT NULL,
    is_size_free boolean DEFAULT false NOT NULL,
    is_weight_free boolean DEFAULT false NOT NULL
);


ALTER TABLE public.cells OWNER TO devuser;

--
-- Name: cells_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.cells_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.cells_id_seq OWNER TO devuser;

--
-- Name: cells_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.cells_id_seq OWNED BY public.cells.id;


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
-- Name: whs; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.whs (
    id integer NOT NULL,
    name character varying(50) DEFAULT ''::character varying
);


ALTER TABLE public.whs OWNER TO devuser;

--
-- Name: whs_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.whs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.whs_id_seq OWNER TO devuser;

--
-- Name: whs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.whs_id_seq OWNED BY public.whs.id;


--
-- Name: zones; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.zones (
    id integer NOT NULL,
    name character varying(50) DEFAULT ''::character varying,
    whs_id integer,
    zone_type smallint
);


ALTER TABLE public.zones OWNER TO devuser;

--
-- Name: zones_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.zones_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.zones_id_seq OWNER TO devuser;

--
-- Name: zones_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.zones_id_seq OWNED BY public.zones.id;


--
-- Name: cells id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.cells ALTER COLUMN id SET DEFAULT nextval('public.cells_id_seq'::regclass);


--
-- Name: whs id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.whs ALTER COLUMN id SET DEFAULT nextval('public.whs_id_seq'::regclass);


--
-- Name: zones id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.zones ALTER COLUMN id SET DEFAULT nextval('public.zones_id_seq'::regclass);


--
-- Data for Name: cells; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.cells (id, name, whs_id, zone_id, passage_id, rack_id, floor, size_length, size_width, size_height, size_volume, size_usefull_volume, size_weight, is_size_free, is_weight_free) FROM stdin;
1	test 1	1	1	2	0	3	0	0	0	0	0	0	f	f
2	invalid storage	99	1	1	0	1	0	0	0	0	0	0	f	f
\.


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
0	2	32	-180
\.


--
-- Data for Name: whs; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.whs (id, name) FROM stdin;
\.


--
-- Data for Name: zones; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.zones (id, name, whs_id, zone_type) FROM stdin;
\.


--
-- Name: cells_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.cells_id_seq', 2, true);


--
-- Name: whs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.whs_id_seq', 1, false);


--
-- Name: zones_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.zones_id_seq', 1, false);


--
-- Name: cells cells_pk; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.cells
    ADD CONSTRAINT cells_pk PRIMARY KEY (id);


--
-- Name: whs whs_pk; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.whs
    ADD CONSTRAINT whs_pk PRIMARY KEY (id);


--
-- Name: zones zones_pk; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.zones
    ADD CONSTRAINT zones_pk PRIMARY KEY (id);


--
-- Name: cells_id_uindex; Type: INDEX; Schema: public; Owner: devuser
--

CREATE UNIQUE INDEX cells_id_uindex ON public.cells USING btree (id);


--
-- Name: cells_whs_id_zone_id_passage_id_floor_uindex; Type: INDEX; Schema: public; Owner: devuser
--

CREATE UNIQUE INDEX cells_whs_id_zone_id_passage_id_floor_uindex ON public.cells USING btree (whs_id, zone_id, passage_id, floor);


--
-- Name: storage1 storage1_cells_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.storage1
    ADD CONSTRAINT storage1_cells_id_fk FOREIGN KEY (cell_id) REFERENCES public.cells(id);


--
-- PostgreSQL database dump complete
--

