CREATE TABLE public.event_txn_id_device_id (
    event_id text NOT NULL,
    room_id text NOT NULL,
    user_id text NOT NULL,
    device_id text NOT NULL,
    txn_id text NOT NULL,
    inserted_ts bigint NOT NULL
);
