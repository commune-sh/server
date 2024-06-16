CREATE TABLE public.per_user_experimental_features (
    user_id text NOT NULL,
    feature text NOT NULL,
    enabled boolean DEFAULT false
);
