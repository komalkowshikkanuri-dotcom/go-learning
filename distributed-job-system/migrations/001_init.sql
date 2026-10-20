
CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying(100),
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE SEQUENCE public.users_id_seq
    AS integer START WITH 1 INCREMENT BY 1
    NO MINVALUE NO MAXVALUE CACHE 1;

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;

ALTER TABLE ONLY public.users
    ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


CREATE TABLE public.jobs (
    id integer NOT NULL,
    user_id integer,
    type character varying(100),
    status character varying(100),
    CONSTRAINT jobs_pkey PRIMARY KEY (id),
    CONSTRAINT jobs_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES public.users(id)
);

CREATE SEQUENCE public.jobs_id_seq
    AS integer START WITH 1 INCREMENT BY 1
    NO MINVALUE NO MAXVALUE CACHE 1;

ALTER SEQUENCE public.jobs_id_seq OWNED BY public.jobs.id;

ALTER TABLE ONLY public.jobs
    ALTER COLUMN id SET DEFAULT nextval('public.jobs_id_seq'::regclass);


CREATE TABLE public.job_attempts (
    id integer NOT NULL,
    job_id integer,
    status character varying(100),
    started_at timestamp with time zone DEFAULT now(),
    finished_at timestamp with time zone,
    CONSTRAINT job_attempts_pkey PRIMARY KEY (id),
    CONSTRAINT job_attempts_job_id_fkey
        FOREIGN KEY (job_id) REFERENCES public.jobs(id)
);

CREATE SEQUENCE public.job_attempts_id_seq
    AS integer START WITH 1 INCREMENT BY 1
    NO MINVALUE NO MAXVALUE CACHE 1;

ALTER SEQUENCE public.job_attempts_id_seq OWNED BY public.job_attempts.id;

ALTER TABLE ONLY public.job_attempts
    ALTER COLUMN id SET DEFAULT nextval('public.job_attempts_id_seq'::regclass);
