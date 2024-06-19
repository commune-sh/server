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
JOIN room_stats_current rsc ON rsc.room_id = r.room_id
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
AND ejj.json::jsonb->'content'->>'public' = 'true'
ORDER BY r.room_id, rsc.joined_members ASC
LIMIT sqlc.narg('limit')::bigint;

-- name: GetPublicSpacesCount :one
SELECT COUNT (DISTINCT r.room_id) as total
FROM rooms r
JOIN current_state_events cse ON r.room_id = cse.room_id
JOIN current_state_events cs ON r.room_id = cs.room_id
LEFT JOIN event_json ej ON ej.event_id = cs.event_id
JOIN current_state_events csp ON r.room_id = csp.room_id
LEFT JOIN event_json ejj ON ejj.event_id = csp.event_id
JOIN room_stats_current rsc ON rsc.room_id = r.room_id
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
SELECT cs.is_parent, rss.*
FROM room_hierarchy rh
JOIN room_stats_state rss 
ON rss.room_id = rh.room_id
JOIN rooms r 
ON r.room_id = rh.room_id 
JOIN current_state_events cse 
ON cse.room_id = r.room_id AND cse.type = 'commune.room.public'
JOIN state_events sev 
ON sev.room_id = r.room_id AND sev.type = 'm.room.create'
JOIN events ev 
ON ev.event_id = sev.event_id
LEFT JOIN LATERAL (
	SELECT EXISTS (SELECT DISTINCT c.type 
	FROM current_state_events c
	WHERE c.room_id = rh.room_id AND c.type = 'm.space.child') as is_parent
) cs ON TRUE
ORDER BY ev.origin_server_ts ASC;

-- name: GetCurrentStateEvents :many
SELECT cse.type as current_state_event, 
    ej.json as event_json, cse.event_id
FROM current_state_events cse
JOIN event_json ej 
ON ej.event_id = cse.event_id
LEFT JOIN current_state_events cs
ON cs.type = 'commune.room.public' AND cs.room_id = cse.state_key
WHERE cse.room_id = $1
AND cse.type != 'm.room.member'
AND 
CASE WHEN cse.type = 'm.space.child' THEN cs.type = 'commune.room.public' 
    ELSE cs.type IS NULL
END;

-- name: GetSpaceChildStateEvents :many
SELECT cse.type as current_state_event, 
    ej.json as event_json, cse.event_id
FROM current_state_events cse
JOIN event_json ej 
ON ej.event_id = cse.event_id
LEFT JOIN current_state_events cs
ON cs.type = 'commune.room.public' AND cs.room_id = cse.state_key
WHERE cse.room_id = $1
AND cse.type = 'm.space.child'
AND 
CASE WHEN cse.type = 'm.space.child' THEN cs.type = 'commune.room.public' 
    ELSE cs.type IS NULL
END;
