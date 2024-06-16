-- name: DoesRoomExist :one
SELECT exists(select 1 from rooms where room_id = $1);

-- name: IsRoomPublic :one
SELECT exists(select 1 from rooms where room_id = $1 and is_public = true);

-- name: IsRoomPubliclyAccessible :one
SELECT EXISTS (
	SELECT r.room_id
	FROM rooms r
	JOIN current_state_events cse ON r.room_id = cse.room_id
	JOIN event_json ej ON ej.event_id = cse.event_id
	WHERE cse.type = 'commune.room.public'
	AND r.room_id = $1
);

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



