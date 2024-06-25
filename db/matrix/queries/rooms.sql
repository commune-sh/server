-- name: DoesRoomExist :one
SELECT exists(select 1 from rooms where room_id = $1);

-- name: IsRoomPublic :one
SELECT exists(select 1 from rooms where room_id = $1 and is_public = true);

-- name: IsRoomPubliclyAccessible :one
SELECT EXISTS (
	SELECT r.room_id
	FROM rooms r
	JOIN room_stats_state rss
	ON rss.room_id = r.room_id
	WHERE r.room_id = $1
	AND rss.join_rules = 'public'
	AND rss.history_visibility = 'world_readable'
	AND r.room_id IN (
		SELECT room_id
		FROM current_state_events 
		WHERE room_id = $1
		AND type = 'commune.room.public'
	)
) is_public;

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

-- name: GetRoomJoinedMembersCount :one
SELECT joined_members 
FROM room_stats_current
WHERE room_id = $1;


-- name: GetRoomState :one
SELECT json_build_object(
		'room_id', rss.room_id,
		'name', rss.name,
        'canonical_alias', rss.canonical_alias,
		'topic', rss.topic, 
		'avatar_url', rss.avatar,
		'join_rule', rss.join_rules,
        'room_type', rss.room_type,
        'guest_can_join', CASE WHEN rss.guest_access = 'can_join' THEN true ELSE false END,
		'world_readable', CASE WHEN rss.history_visibility = 'world_readable' THEN true ELSE false END,
		'num_joined_members', rsc.joined_members
	) as state
FROM room_stats_state rss
JOIN room_stats_current rsc 
ON rss.room_id = rsc.room_id
WHERE rss.room_id = $1
AND rss.join_rules = 'public'
AND rss.history_visibility = 'world_readable'
AND rss.guest_access = 'can_join'
AND rss.room_id IN (
	SELECT room_id
	FROM current_state_events 
	WHERE room_id = $1
	AND type = 'commune.room.public'
);

-- name: GetRoomJoinedMembers :many
SELECT rm.user_id, rm.event_id, rm.display_name, rm.avatar_url, rm.room_id, 
ej.json as event_json
FROM room_memberships rm
JOIN event_json ej
ON ej.event_id = rm.event_id
WHERE rm.room_id = $1
AND rm.membership = 'join'
LIMIT sqlc.narg('limit')::bigint;

