-- "id" bigint DEFAULT nextval('public.users_id_seq'::regclass) NOT NULL,
--  "first_name" character varying(255) NOT NULL,
--  "last_name" character varying(255) NOT NULL,
--  "email" character varying(255) NOT NULL,
--  "password" character varying(255) NOT NULL,
--  "created_at" timestamp without time zone NOT NULL,
--  "updated_at" timestamp without time zone NOT NULL

-- name: CreateUser :one
INSERT INTO users (
  first_name,
  last_name,
  email,
  password,
  created_at,
  updated_at
) VALUES(
$1, $2, $3, $4,$5,$6
) RETURNING *;


-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;
