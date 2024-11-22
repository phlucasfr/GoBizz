package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"auth-service/internal/domain"
	"auth-service/pkg/util"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CompanyRepository struct {
	db      *pgxpool.Pool
	redis   *redis.Client
	queries *Queries
}

func NewCompanyRepository(db *pgxpool.Pool, redis *redis.Client) *CompanyRepository {
	return &CompanyRepository{
		db:      db,
		redis:   redis,
		queries: New(db),
	}
}

func (r *CompanyRepository) Create(ctx context.Context, req domain.CreateCompanyRequest) (*domain.CreateCompanyResponse, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar a transação: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("erro ao hashear a senha: %w", err)
	}

	env := util.GetConfig(".")
	safeCpfCnpj, err := util.Encrypt(req.CPFCNPJ, env.MasterKey)
	if err != nil {
		return nil, fmt.Errorf("erro ao criptografar CPF/CNPJ: %w", err)
	}

	params := CreateCompanyParams{
		Name:           req.Name,
		Email:          req.Email,
		Phone:          req.Phone,
		CpfCnpj:        safeCpfCnpj,
		HashedPassword: string(hashedPassword),
	}

	company, err := r.queries.WithTx(tx).CreateCompany(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar a empresa no banco de dados: %w", err)
	}

	err = util.SendVerificationSms(&company.Phone)
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar SMS de verificação: %w", err)
	}

	return &domain.CreateCompanyResponse{
		ID:   company.ID.Bytes,
		Name: company.Name,
	}, nil
}

func (r *CompanyRepository) VerifyCompanyBySms(ctx context.Context, req *domain.VerifyCompanyBySmsRequest) (bool, error) {
	company, err := r.GetByID(ctx, req.ID)
	if err != nil {
		return false, err
	}

	verified, err := util.CheckVerificationCode(&company.Phone, &req.Code)
	if err != nil {
		return false, err
	}

	if verified {
		_, err = r.queries.ActivateCompany(ctx, company.ID)

		if err != nil {
			return false, err
		}
	}

	return verified, nil
}

func (r *CompanyRepository) GetByID(ctx context.Context, id uuid.UUID) (*Company, error) {
	cacheKey := fmt.Sprintf("company:%s", id.String())
	cachedCompany, err := r.redis.Get(ctx, cacheKey).Result()

	if err == nil && cachedCompany != "" {
		var company Company
		err := json.Unmarshal([]byte(cachedCompany), &company)
		if err != nil {
			return nil, fmt.Errorf("erro ao desserializar dados do cache: %v", err)
		}
		return &company, nil
	}

	pgUUID := pgtype.UUID{
		Bytes: id,
		Valid: true,
	}

	company, err := r.queries.GetCompanyByID(ctx, pgUUID)
	if err != nil {
		return nil, err
	}

	companyJSON, err := json.Marshal(company)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar dados da empresa: %v", err)
	}

	err = r.redis.Set(ctx, cacheKey, companyJSON, 2*time.Minute).Err()
	if err != nil {
		return nil, fmt.Errorf("erro ao armazenar no cache Redis: %v", err)
	}

	return &Company{
		ID:             company.ID,
		Name:           company.Name,
		Email:          company.Email,
		Phone:          company.Phone,
		CpfCnpj:        company.CpfCnpj,
		HashedPassword: company.HashedPassword,
		CreatedAt:      company.CreatedAt,
		UpdatedAt:      company.UpdatedAt,
	}, nil
}
