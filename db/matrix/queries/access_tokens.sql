-- name: IsAccessTokenValid :one
SELECT EXISTS(
    SELECT 1 from access_tokens ac
    WHERE ac.token = $1) AS valid,
    user_id 
FROM access_tokens 
WHERE token = $1;


