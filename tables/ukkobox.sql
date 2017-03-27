--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: ukkobox_download; Type: TABLE; Schema: public; Owner: ukkobox; Tablespace: 
--

CREATE TABLE ukkobox_download (
    up_id integer NOT NULL,
    up_user character varying(300) NOT NULL,
    up_file character varying(300) NOT NULL,
    up_cloud character varying(50) NOT NULL,
    up_create date DEFAULT ('now'::text)::date NOT NULL,
    up_createh time without time zone DEFAULT ('now'::text)::time with time zone NOT NULL,
    up_sent integer DEFAULT 1 NOT NULL
);


ALTER TABLE public.ukkobox_download OWNER TO ukkobox;

--
-- Name: ukkobox_download_up_id_seq; Type: SEQUENCE; Schema: public; Owner: ukkobox
--

CREATE SEQUENCE ukkobox_download_up_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ukkobox_download_up_id_seq OWNER TO ukkobox;

--
-- Name: ukkobox_download_up_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: ukkobox
--

ALTER SEQUENCE ukkobox_download_up_id_seq OWNED BY ukkobox_download.up_id;


--
-- Name: ukkobox_uploads; Type: TABLE; Schema: public; Owner: ukkobox; Tablespace: 
--

CREATE TABLE ukkobox_uploads (
    up_id integer NOT NULL,
    up_user character varying(300) NOT NULL,
    up_file character varying(300) NOT NULL,
    up_cloud character varying(50) NOT NULL,
    up_size character varying(50) NOT NULL,
    up_create date DEFAULT ('now'::text)::date NOT NULL,
    up_createh time without time zone DEFAULT ('now'::text)::time with time zone NOT NULL,
    up_sent integer DEFAULT 1 NOT NULL,
    up_update date DEFAULT ('now'::text)::date NOT NULL,
    up_updateh time without time zone DEFAULT ('now'::text)::time with time zone NOT NULL
);


ALTER TABLE public.ukkobox_uploads OWNER TO ukkobox;

--
-- Name: COLUMN ukkobox_uploads.up_sent; Type: COMMENT; Schema: public; Owner: ukkobox
--

COMMENT ON COLUMN ukkobox_uploads.up_sent IS '1 = local, 2 = sent cloud, 3 = Send error, 4 = in progress, 5 = Removed locally, 6 = Lock to run, progress';


--
-- Name: ukkobox_uploads_up_id_seq; Type: SEQUENCE; Schema: public; Owner: ukkobox
--

CREATE SEQUENCE ukkobox_uploads_up_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ukkobox_uploads_up_id_seq OWNER TO ukkobox;

--
-- Name: ukkobox_uploads_up_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: ukkobox
--

ALTER SEQUENCE ukkobox_uploads_up_id_seq OWNED BY ukkobox_uploads.up_id;


--
-- Name: ukkobox_user; Type: TABLE; Schema: public; Owner: ukkobox; Tablespace: 
--

CREATE TABLE ukkobox_user (
    user_id integer NOT NULL,
    user_name character varying(100) NOT NULL,
    user_email character varying(200) NOT NULL,
    user_pass character varying(100) NOT NULL,
    user_active integer DEFAULT 1 NOT NULL,
    user_token character varying(100) NOT NULL,
    user_date date DEFAULT ('now'::text)::date NOT NULL,
    user_hour time without time zone DEFAULT ('now'::text)::time with time zone NOT NULL
);


ALTER TABLE public.ukkobox_user OWNER TO ukkobox;

--
-- Name: ukkobox_user_user_id_seq; Type: SEQUENCE; Schema: public; Owner: ukkobox
--

CREATE SEQUENCE ukkobox_user_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ukkobox_user_user_id_seq OWNER TO ukkobox;

--
-- Name: ukkobox_user_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: ukkobox
--

ALTER SEQUENCE ukkobox_user_user_id_seq OWNED BY ukkobox_user.user_id;


--
-- Name: up_id; Type: DEFAULT; Schema: public; Owner: ukkobox
--

ALTER TABLE ONLY ukkobox_download ALTER COLUMN up_id SET DEFAULT nextval('ukkobox_download_up_id_seq'::regclass);


--
-- Name: up_id; Type: DEFAULT; Schema: public; Owner: ukkobox
--

ALTER TABLE ONLY ukkobox_uploads ALTER COLUMN up_id SET DEFAULT nextval('ukkobox_uploads_up_id_seq'::regclass);


--
-- Name: user_id; Type: DEFAULT; Schema: public; Owner: ukkobox
--

ALTER TABLE ONLY ukkobox_user ALTER COLUMN user_id SET DEFAULT nextval('ukkobox_user_user_id_seq'::regclass);


--
-- Name: ukkobox_download_pkey; Type: CONSTRAINT; Schema: public; Owner: ukkobox; Tablespace: 
--

ALTER TABLE ONLY ukkobox_download
    ADD CONSTRAINT ukkobox_download_pkey PRIMARY KEY (up_id);


--
-- Name: ukkobox_uploads_pkey; Type: CONSTRAINT; Schema: public; Owner: ukkobox; Tablespace: 
--

ALTER TABLE ONLY ukkobox_uploads
    ADD CONSTRAINT ukkobox_uploads_pkey PRIMARY KEY (up_id);


--
-- Name: ukkobox_uploads_up_user_up_file_key; Type: CONSTRAINT; Schema: public; Owner: ukkobox; Tablespace: 
--

ALTER TABLE ONLY ukkobox_uploads
    ADD CONSTRAINT ukkobox_uploads_up_user_up_file_key UNIQUE (up_user, up_file);


--
-- Name: ukkobox_user_pkey; Type: CONSTRAINT; Schema: public; Owner: ukkobox; Tablespace: 
--

ALTER TABLE ONLY ukkobox_user
    ADD CONSTRAINT ukkobox_user_pkey PRIMARY KEY (user_id);


--
-- Name: ukkobox_user_user_email_key; Type: CONSTRAINT; Schema: public; Owner: ukkobox; Tablespace: 
--

ALTER TABLE ONLY ukkobox_user
    ADD CONSTRAINT ukkobox_user_user_email_key UNIQUE (user_email);


--
-- Name: ukkobox_user_user_token_key; Type: CONSTRAINT; Schema: public; Owner: ukkobox; Tablespace: 
--

ALTER TABLE ONLY ukkobox_user
    ADD CONSTRAINT ukkobox_user_user_token_key UNIQUE (user_token);


--
-- Name: ukkobox_uploads_up_user_fkey; Type: FK CONSTRAINT; Schema: public; Owner: ukkobox
--

ALTER TABLE ONLY ukkobox_uploads
    ADD CONSTRAINT ukkobox_uploads_up_user_fkey FOREIGN KEY (up_user) REFERENCES ukkobox_user(user_token) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

