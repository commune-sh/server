-- name: IsRoomPublic :one
SELECT exists(select 1 from rooms where room_id = $1 and is_public = true);
