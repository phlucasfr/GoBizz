package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	IpHash string `json:"ip_hash"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(ConfigInstance.MasterKey)

// ComputeIpHash generates a SHA-256 hash based on the provided IP address and user agent string.
// The IP address and user agent are concatenated with a "|" separator before hashing.
// The resulting hash is returned as a hexadecimal-encoded string.
//
// Parameters:
//   - ip: The IP address as a string.
//   - ua: The user agent string.
//
// Returns:
//
//	A hexadecimal-encoded string representing the SHA-256 hash of the concatenated input.
func ComputeIpHash(ip, ua string) string {
	h := sha256.Sum256([]byte(ip + "|" + ua))
	return hex.EncodeToString(h[:])
}

// GenerateJWT generates a JSON Web Token (JWT) for a user with the provided details.
// It includes custom claims such as user ID, email, and a hashed representation of the user's IP and User-Agent.
// The token is signed using the HS256 algorithm.
//
// Parameters:
//   - userID: The unique identifier of the user.
//   - email: The email address of the user.
//   - ip: The IP address of the user.
//   - userAgent: The User-Agent string of the user's device.
//
// Returns:
//   - string: The signed JWT as a string.
//   - error: An error if the token generation or signing fails.
func GenerateJWT(userID, email, ip, userAgent string) (string, error) {
	ipHash := ComputeIpHash(ip, userAgent)

	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		IpHash: ipHash,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth-service",
			Audience:  jwt.ClaimStrings{"auth-service"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if token.Method != jwt.SigningMethodHS256 {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
