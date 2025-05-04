package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTClaims(t *testing.T) {
	masterKey := ConfigInstance.MasterKey
	t.Run("Properly verifies expiration", func(t *testing.T) {
		ip := "127.0.0.1"
		userAgent := "Mozilla/5.0"
		token, _ := GenerateJWT("testUser", "testEmail@example.com", ip, userAgent)

		parsedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return masterKey, nil
		})

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		require.True(t, ok, "Expected claims to be of type jwt.MapClaims")

		exp, ok := claims["exp"].(float64)
		require.True(t, ok, "Expected expiration claim to be a float64")
		expTime := time.Unix(int64(exp), 0)
		require.True(t, expTime.After(time.Now()), "Expected token expiration time to be in the future")
	})
}

func TestValidateJWT(t *testing.T) {
	masterKey := ConfigInstance.MasterKey
	t.Run("Malformed token", func(t *testing.T) {
		_, err := ValidateJWT("token_invalido.formato")
		require.Error(t, err)
	})

	t.Run("Expired token", func(t *testing.T) {
		claims := JWTClaims{
			UserID: "expiredUser",
			Email:  "test@example.com",
			IpHash: "expiredIpHash",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(masterKey))

		_, err := ValidateJWT(tokenString)
		require.Error(t, err)
		require.Contains(t, err.Error(), "token is expired")
	})

	t.Run("Token without subject", func(t *testing.T) {
		claims := jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(masterKey)

		_, err := ValidateJWT(tokenString)
		require.Error(t, err)
	})

	t.Run("Valid token", func(t *testing.T) {
		claims := JWTClaims{
			UserID: "validUser",
			Email:  "valid@example.com",
			IpHash: ComputeIpHash("127.0.0.1", "Mozilla/5.0"),
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(masterKey))

		validatedClaims, err := ValidateJWT(tokenString)
		require.NoError(t, err)
		require.NotNil(t, validatedClaims)
		require.Equal(t, claims.UserID, validatedClaims.UserID)
		require.Equal(t, claims.Email, validatedClaims.Email)
		require.Equal(t, claims.IpHash, validatedClaims.IpHash)
	})
}
