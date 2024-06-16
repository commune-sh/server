CREATE TABLE public.scheduled_tasks (
    id text NOT NULL,
    action text NOT NULL,
    status text NOT NULL,
    "timestamp" bigint NOT NULL,
    resource_id text,
    params text,
    result text,
    error text
);
