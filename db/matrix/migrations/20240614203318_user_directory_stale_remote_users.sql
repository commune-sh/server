CREATE TABLE public.user_directory_stale_remote_users (
    user_id text NOT NULL,
    user_server_name text NOT NULL,
    next_try_at_ts bigint NOT NULL,
    retry_counter integer NOT NULL
);
