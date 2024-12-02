-- name: CreateCompany :one
INSERT INTO company (name, email, phone, cpf_cnpj, is_active, hashed_password)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetCompanyByID :one
SELECT * FROM company WHERE id = $1;

-- name: GetCompanyByEmail :one
SELECT * FROM company WHERE email = $1;

-- name: HasActiveCompany :one
SELECT EXISTS (
  SELECT 1 FROM company WHERE (email = $1 OR phone = $2) AND is_active = true
);

-- name: ListCompanies :many
SELECT * FROM company ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: UpdatePasswordByEmail :execrows
UPDATE company
SET
  hashed_password = $2,
  updated_at = NOW()
WHERE email = $1;

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

-- name: DeleteCompany :execrows
DELETE FROM company WHERE id = $1;