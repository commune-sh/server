-- name: SearchRoomMessages :many
SELECT ts_rank_cd(vector, websearch_to_tsquery('english', $1)) AS rank,
es.room_id, es.event_id, ej.json
FROM event_search es
JOIN room_stats_state rss
ON rss.room_id = es.room_id
JOIN event_json ej
ON ej.event_id = es.event_id
WHERE es.vector @@  websearch_to_tsquery('english', $1)
AND rss.join_rules = 'public'
AND rss.history_visibility = 'world_readable'
AND es.room_id = $1
AND es.room_id IN (
	SELECT room_id
	FROM current_state_events 
	WHERE room_id = $1
	AND type = 'commune.room.public'
);

-- name: SearchPublicRoomMessages :many
SELECT ts_rank_cd(vector, websearch_to_tsquery('english', $1)) AS rank,
es.room_id, es.event_id, ej.json
FROM event_search es
JOIN room_stats_state rss
ON rss.room_id = es.room_id
JOIN event_json ej
ON ej.event_id = es.event_id
WHERE es.vector @@  websearch_to_tsquery('english', $1)
AND rss.join_rules = 'public'
AND rss.history_visibility = 'world_readable'
AND es.room_id IN (
	SELECT room_id
	FROM current_state_events 
	WHERE type = 'commune.room.public'
);
