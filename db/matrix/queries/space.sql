-- name: GetSpaceChildren :many
SELECT cse.state_key as room_id
FROM current_state_events cse
WHERE cse.room_id = $1
AND cse.type = 'm.space.child';

-- name: GetSpaceChildParent :one
SELECT cse.state_key as room_id
FROM current_state_events cse
WHERE cse.room_id = $1
AND cse.type = 'm.space.parent';

-- name: GetSpaceChildParent :one
