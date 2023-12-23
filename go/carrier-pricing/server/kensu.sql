--
-- PostgreSQL database dump
--

-- Dumped from database version 15.2
-- Dumped by pg_dump version 15.2

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
-- Name: carriers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.carriers (
    id integer NOT NULL,
    name character varying(86) NOT NULL,
    price integer NOT NULL,
    markup integer NOT NULL,
    "time" integer NOT NULL,
    vehicles text[] NOT NULL
);


ALTER TABLE public.carriers OWNER TO postgres;

--
-- Name: carriers_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.carriers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.carriers_id_seq OWNER TO postgres;

--
-- Name: carriers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.carriers_id_seq OWNED BY public.carriers.id;


--
-- Name: carriers id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.carriers ALTER COLUMN id SET DEFAULT nextval('public.carriers_id_seq'::regclass);


--
-- Data for Name: carriers; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.carriers (id, name, price, markup, "time", vehicles) FROM stdin;
1	CollectTimes	50	20	1	{parcel_car,small_van,large_van}
2	CollectTimes	50	10	3	{parcel_car,motorbike}
3	CollectTimes	50	5	5	{bicycle,motorbike}
4	RoyalPackages	30	5	3	{bicycle,motorbike,parcel_car}
5	RoyalPackages	30	50	1	{bicycle,motorbike,parcel_car,small_van,large_van}
6	Hercules	25	10	5	{motorbike,parcel_car,small_van}
7	Hercules	25	0	10	{large_van,parcel_car,bicycle}
8	OOPS	0	0	12	{large_van,parcel_car,bicycle}
9	OOPS	0	0	7	{large_van,motorbike,parcel_car,small_van}
\.


--
-- Name: carriers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.carriers_id_seq', 334, true);


--
-- Name: carriers carriers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.carriers
    ADD CONSTRAINT carriers_pkey PRIMARY KEY (id);


--
-- Name: carriers uni; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.carriers
    ADD CONSTRAINT uni UNIQUE (name, vehicles);


--
-- PostgreSQL database dump complete
--

