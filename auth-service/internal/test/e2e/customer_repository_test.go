package e2e

import (
	"auth-service/internal/domain"
	"auth-service/internal/test"
	"auth-service/utils"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateCompany(t *testing.T) {
	rep, err := test.SetupTestContainers()
	require.NoError(t, err, "failed to set up test containers")

	ctx := context.Background()
	params := domain.CreateCustomerRequest{
		Name:     utils.RandomString(16),
		Email:    "test@gmail.com",
		Phone:    "47999999999",
		CPFCNPJ:  "99999999999",
		Password: utils.RandomString(8),
	}

	company, err := rep.CustomerRepository.Create(ctx, params)
	require.NoError(t, err, "failed to create company")
	require.NotEmpty(t, company.ID, "expected company ID to be non-empty")
	require.Equal(t, params.Name, company.Name, "company name should match the provided name")
}
