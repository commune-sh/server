-- name: GetStateEventsOfType :many
SELECT cse.event_id, cse.state_key
FROM current_state_events cse
JOIN room_stats_state rss 
ON rss.room_id = cse.room_id
WHERE cse.room_id = $1
AND cse.type = 'm.room.member'
AND rss.join_rules = 'public'
AND rss.history_visibility = 'world_readable'
AND rss.guest_access = 'can_join'
AND rss.room_id IN (
	SELECT room_id
	FROM current_state_events 
	WHERE room_id = $1
	AND type = 'commune.room.public'
);



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

