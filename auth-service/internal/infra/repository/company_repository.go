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
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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

	err = r.redis.Set(ctx, cacheKey, companyJSON, 30*time.Minute).Err()
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

func (r *CompanyRepository) VerifyCompanyBySms(ctx context.Context, req *domain.VerifyCompanyBySmsRequest) error {
	company, err := r.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	verified, err := util.CheckVerificationCode(&company.Phone, &req.Code)
	if err != nil {
		return err
	}

	if verified {
		err := r.persistCompanyInDB(ctx, company)
		if err != nil {
			return err
		}

		cacheKey := fmt.Sprintf("company:%x", company.ID.String())
		err = r.redis.Del(ctx, cacheKey).Err()
		if err != nil {
			return fmt.Errorf("erro ao remover empresa do cache Redis após persistência: %v", err)
		}
	}

	return nil
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

	pgID := pgtype.UUID{
		Bytes: id,
		Valid: true,
	}

	company, err := r.queries.GetCompanyByID(ctx, pgID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar empresa no banco de dados: %v", err)
	}

	err = r.redis.Set(ctx, cacheKey, company, 2*time.Minute).Err()
	if err != nil {
		return nil, fmt.Errorf("erro ao salvar dados no cache: %v", err)
	}

	return &domain.Company{
		ID:             company.ID.Bytes,
		Name:           company.Name,
		Email:          company.Email,
		Phone:          company.Phone,
		CpfCnpj:        company.CpfCnpj,
		IsActive:       company.IsActive,
		UpdatedAt:      company.UpdatedAt,
		CreatedAt:      company.CreatedAt,
		HashedPassword: company.HashedPassword,
	}, nil
}

func (r *CompanyRepository) SendRecoveryEmail(ctx context.Context, token string, email string) error {
	cacheKey := fmt.Sprintf("password-recovery:%s", token)
	cacheValue := map[string]string{"email": email}

	cacheValueJSON, err := json.Marshal(cacheValue)
	if err != nil {
		return fmt.Errorf("erro ao serializar dados para o Redis: %v", err)
	}

	err = r.redis.Set(ctx, cacheKey, cacheValueJSON, 10*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("erro ao armazenar token no Redis: %v", err)
	}

	env := util.GetConfig(".")
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", env.FrontendSource, token)

	from := mail.NewEmail("GoBizz", "gobizz.comercial@gmail.com")
	subject := "Recuperação de Senha"
	to := mail.NewEmail("Usuário", email)
	plainTextContent := fmt.Sprintf("Clique no link para redefinir sua senha: %s", resetLink)
	htmlContent := fmt.Sprintf(`<p>Clique no link para redefinir sua senha: <a href="%s">%s</a></p>`, resetLink, resetLink)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(env.SendGridApiKey)
	_, err = client.Send(message)
	if err != nil {
		return fmt.Errorf("erro ao enviar email: %v", err)
	}

	return nil
}

func (r *CompanyRepository) ValidateResetToken(ctx context.Context, token string) (string, error) {
	cacheKey := fmt.Sprintf("password-recovery:%s", token)

	cacheValueJSON, err := r.redis.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("token inválido ou expirado")
	} else if err != nil {
		return "", fmt.Errorf("erro ao acessar o Redis: %v", err)
	}

	var cacheValue map[string]string
	err = json.Unmarshal([]byte(cacheValueJSON), &cacheValue)
	if err != nil {
		return "", fmt.Errorf("erro ao desserializar dados do Redis: %v", err)
	}

	email, ok := cacheValue["email"]
	if !ok {
		return "", fmt.Errorf("email não encontrado para o token fornecido")
	}

	return email, nil
}

func (r *CompanyRepository) UpdatePasswordByEmail(ctx context.Context, token string, email string, hashedPassword string) error {
	_, err := r.queries.UpdatePasswordByEmail(ctx, UpdatePasswordByEmailParams{
		Email:          email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return fmt.Errorf("erro ao atualizar senha no banco de dados: %v", err)
	}

	tokenKey := fmt.Sprintf("password-recovery:%s", token)
	err = r.redis.Del(ctx, tokenKey).Err()
	if err != nil {
		return fmt.Errorf("erro ao remover token do Redis: %v", err)
	}

	return nil
}

func (r *CompanyRepository) GetByEmail(ctx context.Context, email string) (*domain.Company, error) {
	company, err := r.queries.GetCompanyByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	return &domain.Company{
		ID:             company.ID.Bytes,
		Name:           company.Name,
		Email:          company.Email,
		Phone:          company.Phone,
		CpfCnpj:        company.CpfCnpj,
		IsActive:       company.IsActive,
		UpdatedAt:      company.UpdatedAt,
		CreatedAt:      company.CreatedAt,
		HashedPassword: company.HashedPassword,
	}, nil
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
