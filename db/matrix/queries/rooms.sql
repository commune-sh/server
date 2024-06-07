-- name: IsRoomPublic :one
SELECT exists(select 1 from rooms where room_id = $1 and is_public = true);

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

-- name: GetRoomCurrentStateEvents :many
SELECT ej.json::jsonb->>'content' as content
FROM current_state_events cse
JOIN event_json ej ON ej.event_id = cse.event_id
WHERE ej.room_id = $1;
