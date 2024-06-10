package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
)

// VerifyJWT verifies the JWT token and returns the claims
func VerifyJWT(tokenString string) (jwt.MapClaims, error) {
	var secretKey = []byte("c2t1bGlkc3RyaW5nZm9ya2V5cGVycHJvamVjdGJldHdlZW5wcm9qZWN0")
	// sample token string taken from the New example

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(token.Claims)

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("failed to extract claims")
		return nil, fmt.Errorf("failed to extract claims")
	}

	return claims, nil

}
