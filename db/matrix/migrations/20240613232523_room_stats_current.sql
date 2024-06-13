-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS room_stats_current (
    room_id text NOT NULL,
    current_state_events int NOT NULL,
    joined_members int NOT NULL,
    invited_members int NOT NULL,
    left_members int NOT NULL,
    banned_members int NOT NULL,
    local_users_in_room int NOT NULL,
    completed_delta_stream_id bigint NOT NULL,
    knocked_members int
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE redactions;
-- +goose StatementEnd
