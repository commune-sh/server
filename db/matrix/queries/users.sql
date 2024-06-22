-- name: IsUserAdmin :one
SELECT 
    CASE WHEN u.admin = 1 AND u.suspended IS false THEN TRUE 
        ELSE FALSE 
    END AS admin
FROM users u 
WHERE u.name = $1;
