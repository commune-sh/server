-- name: GetPublicSpacesUnsafe :many
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

-- name: GetPublicSpaces :many
SELECT DISTINCT ON (r.room_id) r.room_id
FROM rooms r
JOIN current_state_events cse ON r.room_id = cse.room_id
JOIN current_state_events cs ON r.room_id = cs.room_id
LEFT JOIN event_json ej ON ej.event_id = cs.event_id
JOIN current_state_events csp ON r.room_id = csp.room_id
LEFT JOIN event_json ejj ON ejj.event_id = csp.event_id
WHERE cse.type = 'm.space.child'
  AND r.room_id NOT IN (
    SELECT room_id
    FROM current_state_events
    WHERE type = 'm.space.parent'
  )
AND r.is_public is true
AND cs.type = 'm.room.history_visibility'
AND ej.json::jsonb->'content'->>'history_visibility' = 'world_readable'
AND csp.type = 'commune.room.public'
AND ejj.json::jsonb->'content'->>'public' = 'true';

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
