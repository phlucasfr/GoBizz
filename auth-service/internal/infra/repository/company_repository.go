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
	hasActiveCompany, err := r.queries.HasActiveCompany(ctx, HasActiveCompanyParams{
		Email: req.Email,
		Phone: req.Phone,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar dados da empresa: %v", err)
	}

	if hasActiveCompany {
		return nil, fmt.Errorf("já existe uma empresa ativa com os dados informados (email ou telefone)")
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	env := util.GetConfig(".")
	encryptedCpfCnpf, err := util.Encrypt(req.CPFCNPJ, env.MasterKey)
	if err != nil {
		return nil, err
	}
	company := domain.Company{
		ID:             uuid.New(),
		Name:           req.Name,
		Email:          req.Email,
		Phone:          req.Phone,
		CpfCnpj:        encryptedCpfCnpf,
		HashedPassword: hashedPassword,
	}

	companyJSON, err := json.Marshal(company)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar os dados da empresa: %v", err)
	}

	cacheKey := fmt.Sprintf("company:%x", company.ID.String())

	err = r.redis.Set(ctx, cacheKey, companyJSON, 24*time.Hour).Err()
	if err != nil {
		return nil, fmt.Errorf("erro ao armazenar empresa no cache Redis: %v", err)
	}

	err = util.SendVerificationSms(&company.Phone)
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar SMS de verificação: %w", err)
	}

	return &domain.CreateCompanyResponse{
		ID:   company.ID,
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
		err := r.persistCompanyInDB(ctx, company)
		if err != nil {
			return false, err
		}

		cacheKey := fmt.Sprintf("company:%x", company.ID.String())
		err = r.redis.Del(ctx, cacheKey).Err()
		if err != nil {
			return false, fmt.Errorf("erro ao remover empresa do cache Redis após persistência: %v", err)
		}
	}

	return verified, nil
}

func (r *CompanyRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	cacheKey := fmt.Sprintf("company:%x", id.String())
	cachedCompany, err := r.redis.Get(ctx, cacheKey).Result()

	if err == nil && cachedCompany != "" {
		var company domain.Company
		err := json.Unmarshal([]byte(cachedCompany), &company)
		if err != nil {
			return nil, fmt.Errorf("erro ao desserializar dados do cache: %v", err)
		}
		return &company, nil
	}

	return nil, fmt.Errorf("empresa não encontrada no cache ou banco de dados")
}

func (r *CompanyRepository) persistCompanyInDB(ctx context.Context, company *domain.Company) error {
	params := CreateCompanyParams{
		Name:           company.Name,
		Email:          company.Email,
		Phone:          company.Phone,
		CpfCnpj:        company.CpfCnpj,
		IsActive:       true,
		HashedPassword: company.HashedPassword,
	}

	_, err := r.queries.CreateCompany(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao persistir empresa no banco de dados: %w", err)
	}

	cacheKey := fmt.Sprintf("company:%x", company.ID)
	err = r.redis.Del(ctx, cacheKey).Err()
	if err != nil {
		return fmt.Errorf("erro ao remover empresa do cache Redis após persistência: %v", err)
	}

	return nil
}
