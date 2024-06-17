-- name: GetRoomMessages :many
SELECT ej.event_id, ej.json 
FROM event_json ej
JOIN events ON events.event_id = ej.event_id
WHERE ej.room_id = $1
ORDER BY CASE
    WHEN @order_by::text = 'ASC' THEN events.origin_server_ts 
END ASC, CASE 
    WHEN @order_by::text = 'DESC' THEN events.origin_server_ts 
END DESC, CASE
    WHEN @order_by::text = '' THEN events.origin_server_ts 
END DESC
LIMIT sqlc.narg('limit')::bigint;

