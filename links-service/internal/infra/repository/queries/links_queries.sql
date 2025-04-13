-- name: CreateLink :one
INSERT INTO links (
    original_url,
    short_url,
    custom_slug,
    customer_id,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetLinkByShortURL :one
SELECT * FROM links
WHERE short_url = $1 LIMIT 1;

-- name: GetLinkByCustomSlug :one
SELECT * FROM links
WHERE custom_slug = $1 LIMIT 1;

-- name: GetLinksByCustomer :many
SELECT 
    l.*,
    COUNT(*) OVER() as total_count
FROM links l
WHERE l.customer_id = $1
    AND ($4::text IS NULL OR 
        LOWER(l.original_url) LIKE '%' || LOWER($4) || '%' OR
        (l.custom_slug IS NOT NULL AND LOWER(l.custom_slug) LIKE '%' || LOWER($4) || '%')
    )
    AND ($5::text IS NULL OR 
        CASE 
            WHEN $5 = 'active' THEN (l.expires_at IS NULL OR l.expires_at > NOW())
            WHEN $5 = 'expired' THEN (l.expires_at IS NOT NULL AND l.expires_at <= NOW())
            ELSE true
        END
    )
    AND ($6::text IS NULL OR 
        CASE 
            WHEN $6 = 'custom' THEN l.custom_slug IS NOT NULL
            WHEN $6 = 'auto' THEN l.custom_slug IS NULL
            ELSE true
        END
    )
ORDER BY 
    CASE $7::text
        WHEN 'clicks' THEN l.clicks::text
        WHEN 'created_at' THEN l.created_at::text
        WHEN 'original_url' THEN l.original_url
        WHEN 'expiration_date' THEN 
            CASE 
                WHEN l.expires_at IS NULL THEN '9999-12-31'::text
                ELSE l.expires_at::text
            END
        ELSE l.created_at::text
    END
    DESC
LIMIT NULLIF($2, 0)
OFFSET COALESCE($3, 0);

-- name: CountLinksByCustomer :one
SELECT COUNT(*) as total
FROM links l
WHERE l.customer_id = $1;

-- name: UpdateLinkClicks :one
UPDATE links
SET clicks = clicks + 1
WHERE id = $1
RETURNING *;

-- name: UpdateLink :one
UPDATE links
SET 
    original_url = COALESCE($1, original_url),
    short_url = COALESCE($2, short_url),
    custom_slug = COALESCE($3, custom_slug),
    expires_at = $4
WHERE id = $5
RETURNING id, original_url, short_url, custom_slug, clicks, created_at, updated_at, expires_at, customer_id;

-- name: DeleteLink :exec
DELETE FROM links
WHERE id = $1;

-- name: GetLinkByID :one
SELECT * FROM links
WHERE id = $1 LIMIT 1;

-- name: CheckShortURLExists :one
SELECT EXISTS (
    SELECT 1 FROM links
    WHERE short_url = $1
);

-- name: CheckCustomSlugExists :one
SELECT EXISTS (
    SELECT 1 FROM links
    WHERE custom_slug = $1
);

-- name: GetExpiredLinks :many
SELECT * FROM links
WHERE expires_at IS NOT NULL 
AND expires_at < NOW();

-- name: DeleteExpiredLinks :exec
DELETE FROM links
WHERE expires_at IS NOT NULL 
AND expires_at < NOW();

CREATE INDEX idx_links_customer_id ON links(customer_id);
CREATE INDEX idx_links_created_at ON links(created_at);
CREATE INDEX idx_links_expires_at ON links(expires_at);





