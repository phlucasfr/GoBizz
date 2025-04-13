-- name: CreateCustomer :one
INSERT INTO customer (name, email, phone, cpf_cnpj, is_active, hashed_password)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetCustomerByID :one
SELECT * FROM customer WHERE id = $1;

-- name: GetCustomerByEmail :one
SELECT * FROM customer WHERE email = $1;

-- name: HasActiveCustomer :one
SELECT EXISTS (
  SELECT 1 FROM customer WHERE (email = $1 OR phone = $2) AND is_active = true
);

-- name: ListCompanies :many
SELECT * FROM customer ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ActivateCustomerByEmail :one
UPDATE customer
SET
  is_active = true,
  updated_at = NOW()
WHERE email = $1
RETURNING *;

-- name: UpdatePasswordByEmail :execrows
UPDATE customer
SET
  hashed_password = $2,
  updated_at = NOW()
WHERE email = $1;

-- name: UpdateCustomer :one
UPDATE customer
SET
  name = COALESCE($2, name),
  email = COALESCE($3, email),
  phone = COALESCE($4, phone),
  cpf_cnpj = COALESCE($5, cpf_cnpj),
  is_active = COALESCE($6, is_active),
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ActivateCustomer :one
UPDATE customer
SET
  is_active = true,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCustomer :execrows
DELETE FROM customer WHERE id = $1;