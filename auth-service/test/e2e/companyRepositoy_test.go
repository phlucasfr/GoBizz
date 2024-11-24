package test

import (
	"context"
	"testing"

	"auth-service/internal/domain"
	"auth-service/pkg/util"
	"auth-service/test"

	"github.com/stretchr/testify/require"
)

func TestCreateCompany(t *testing.T) {
	repositories, err := test.SetupTestContainers(t)
	require.NoError(t, err)

	ctx := context.Background()
	params := domain.CreateCompanyRequest{
		Name:     util.RandomString(16),
		Email:    "test@gmail.com",
		Phone:    "47999999999",
		CPFCNPJ:  "99999999999",
		Password: util.RandomString(8),
	}

	company, err := repositories.CompanyRepository.Create(ctx, params)
	require.NoError(t, err)
	require.NotEmpty(t, company.ID)
	require.Equal(t, company.Name, params.Name)
}
