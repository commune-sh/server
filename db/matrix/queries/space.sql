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

-- name: GetRoomHierarchy :many
WITH RECURSIVE room_hierarchy AS (
    SELECT DISTINCT cse.room_id
    FROM current_state_events cse
    WHERE cse.room_id = $1
    AND cse.type = 'm.space.child'
  
    UNION ALL
  
    SELECT DISTINCT c.state_key
    FROM current_state_events c
    INNER JOIN room_hierarchy rh 
    ON c.room_id = rh.room_id
    WHERE c.type = 'm.space.child'
)
SELECT * FROM room_hierarchy;

-- name: GetCurrentStateEvents :many
SELECT cse.type current_state_event, 
    ej.json as event_json, cse.event_id
FROM current_state_events cse
JOIN event_json ej 
ON ej.event_id = cse.event_id
WHERE cse.room_id = $1;
