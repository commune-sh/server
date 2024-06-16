CREATE TABLE public.profiles (
    user_id text NOT NULL,
    displayname text,
    avatar_url text,
    full_user_id text,
    CONSTRAINT full_user_id_not_null CHECK ((full_user_id IS NOT NULL))
);
