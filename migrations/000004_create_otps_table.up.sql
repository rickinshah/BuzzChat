CREATE TABLE IF NOT EXISTS otps (
    user_pid bigint PRIMARY KEY REFERENCES users ON DELETE CASCADE,
    created_at timestamp(0) with time zone NOT NULL,
    secret_key text NOT NULL,
    expiry timestamp(0) with time zone NOT NULL
);
