-- name: DoesRoomExist :one
SELECT exists(select 1 from rooms where room_id = $1);

-- name: IsRoomPublic :one
SELECT exists(select 1 from rooms where room_id = $1 and is_public = true);

-- name: GetRoomIDFromAlias :one
SELECT rooms.room_id 
FROM rooms JOIN room_aliases 
ON rooms.room_id = room_aliases.room_id 
WHERE room_aliases.room_alias = $1;

-- name: GetRoomAliasFromID :one
SELECT room_aliases.room_alias
FROM rooms JOIN room_aliases
ON rooms.room_id = room_aliases.room_id
WHERE rooms.room_id = $1;

-- name: GetRoomHistoryVisibility :one
SELECT ej.json::jsonb->'content'->>'history_visibility' AS visibility 
FROM current_state_events cse JOIN event_json ej ON ej.event_id = cse.event_id 
WHERE cse.room_id = $1
AND cse.type = 'm.room.history_visibility';

-- name: GetRoomJoinRules :one
SELECT ej.json::jsonb->'content'->>'join_rule' AS visibility 
FROM current_state_events cse JOIN event_json ej ON ej.event_id = cse.event_id 
WHERE cse.room_id = $1
AND cse.type = 'm.room.join_rules';

-- name: GetRoomGuestAccess :one
SELECT ej.json::jsonb->'content'->>'guest_access' AS visibility 
FROM current_state_events cse JOIN event_json ej ON ej.event_id = cse.event_id 
WHERE cse.room_id = $1
AND cse.type = 'm.room.guest_access';

-- name: GetRoomJoinedMembers :one
SELECT joined_members 
FROM room_stats_current
WHERE room_id = $1;


-- name: GetPublicRooms :many
SELECT DISTINCT ON (r.room_id) r.room_id
FROM rooms r
JOIN current_state_events cse ON r.room_id = cse.room_id
JOIN current_state_events cs ON r.room_id = cs.room_id
LEFT JOIN event_json ej ON ej.event_id = cs.event_id
WHERE cse.type = 'm.space.child'
  AND r.room_id NOT IN (
    SELECT room_id
    FROM current_state_events
    WHERE type = 'm.space.parent'
  )
AND r.is_public is true
AND cs.type = 'm.room.history_visibility'
AND ej.json::jsonb->'content'->>'history_visibility' = 'world_readable';

