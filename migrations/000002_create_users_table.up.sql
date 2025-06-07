CREATE TABLE IF NOT EXISTS users (
    user_pid bigint PRIMARY KEY DEFAULT next_id (),
    username citext UNIQUE NOT NULL,
    email citext UNIQUE NOT NULL,
    name citext,
    password_hash bytea NOT NULL,
    bio text,
    activated bool NOT NULL DEFAULT FALSE,
    profile_pic text,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT now(),
    version int NOT NULL DEFAULT 1
);

