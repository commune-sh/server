CREATE TABLE public.room_forgetter_stream_pos (
    lock character(1) DEFAULT 'X'::bpchar NOT NULL,
    stream_id bigint NOT NULL,
    CONSTRAINT room_forgetter_stream_pos_lock_check CHECK ((lock = 'X'::bpchar))
);
