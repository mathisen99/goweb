SELECT
    id,
    username,
    email,
    created_at
FROM
    user
ORDER BY
    username DESC
LIMIT 10;
