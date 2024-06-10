package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Claims Example claims struct
type Claims struct {
	IssueBy      string `json:"issue_by"`
	PermittedFor string `json:"permitted_for"`
	IdentifiedBy string `json:"identified_by"`
	jwt.StandardClaims
}

// Secret key for signing and validating the token
var secretKey = []byte("c2t1bGlkc3RyaW5nZm9ya2V5cGVycHJvamVjdGJldHdlZW5wcm9qZWN0")

func generateJWT(IssueBy string, PermittedFor string, expirationTime time.Time) (string, error) {
	// Create the claims
	claims := &Claims{
		IssueBy:      IssueBy,
		PermittedFor: PermittedFor,
		IdentifiedBy: "12304masdkn",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func JWTGenerator() (string, error) {
	// Example usage
	expirationTime := time.Now().Add(24 * time.Hour)
	token, err := generateJWT("white_tiger", "tiger", expirationTime)
	if err != nil {
		fmt.Println("Error generating token:", err)
	}

	fmt.Println("Generated Token:", token)
	return token, err
}
