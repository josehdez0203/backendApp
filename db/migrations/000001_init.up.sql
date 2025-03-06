
--
-- PostgreSQL database dump
--


CREATE SEQUENCE "public"."users_id_seq"
   START WITH 1
   INCREMENT BY 1
   MINVALUE 1
   MAXVALUE 9223372036854775807
   CACHE 1;

CREATE TABLE "public"."users"(
   "id" bigint DEFAULT nextval('public.users_id_seq'::regclass) NOT NULL,
   "first_name" character varying(255) NOT NULL,
   "last_name" character varying(255) NOT NULL,
   "email" character varying(255) NOT NULL,
   "password" character varying(255) NOT NULL,
   "created_at" timestamp without time zone NOT NULL,
   "updated_at" timestamp without time zone NOT NULL
);

CREATE UNIQUE INDEX users_pkey ON public.users USING btree (id);
CREATE UNIQUE INDEX email ON public.users USING btree (email);
--
-- Data for Name: users
--
insert into
  users (
    first_name,
    last_name,
    email,
    password,
    created_at,
    updated_at
  )
values  (
    'Admin',
    'User',
    'admin@example.com',
    '$2a$14$wVsaPvJnJJsomWArouWCtusem6S/.Gauq/GjOIEHpyh2DAMmso1wy',
    '2022-09-23 00:00:00',
    '2022-09-23 00:00:00'
  );

--
-- PostgreSQL database dump complete
--
