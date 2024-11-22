-- name: CreateCompany :one
INSERT INTO company (name, email, phone, cpf_cnpj, is_active, hashed_password)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetCompanyByID :one
SELECT * FROM company WHERE id = $1;

-- name: ListCompanies :many
SELECT * FROM company ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: UpdateCompany :one
UPDATE company 
SET 
    name = COALESCE($2, name),
    email = COALESCE($3, email),
    phone = COALESCE($4, phone),
    cpf_cnpj = COALESCE($5, cpf_cnpj),
    is_active = COALESCE($6, is_active),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ActivateCompany :one
UPDATE company 
SET 
    is_active = true,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCompany :exec
DELETE FROM company WHERE id = $1;