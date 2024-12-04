package domain

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateCompanyRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	CPFCNPJ  string `json:"cpf_cnpj" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type CreateCompanyResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type PasswordRecoveryRequest struct {
	Email string `json:"email" validate:"required"`
}

type ActivateCompanyByEmailRequest struct {
	Email string `json:"email" validate:"required"`
}

type VerifyCompanyByEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

type Company struct {
	ID             uuid.UUID          `json:"id"`
	Name           string             `json:"name"`
	Email          string             `json:"email"`
	Phone          string             `json:"phone"`
	CpfCnpj        string             `json:"cpf_cnpj"`
	IsActive       bool               `json:"is_active"`
	UpdatedAt      pgtype.Timestamptz `json:"updated_at"`
	CreatedAt      pgtype.Timestamptz `json:"created_at"`
	HashedPassword string             `json:"hashed_password"`
}
