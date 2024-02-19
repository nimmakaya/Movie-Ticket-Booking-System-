package utils

import (
	"backend/models"
	"github.com/dgrijalva/jwt-go"
)

// Helper function to generate JWT
func GenerateJWT(claims models.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": claims.Email,
		// Add any other claims you want to include
	})
	return token.SignedString([]byte("XmjiwGSmY+X/uAUYZxAxVz2NxAcg8EJZsjSCIyUkpNg="))
}
