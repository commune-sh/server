CREATE TABLE public.timeline_gaps (
    room_id text NOT NULL,
    instance_name text NOT NULL,
    stream_ordering bigint NOT NULL
);
