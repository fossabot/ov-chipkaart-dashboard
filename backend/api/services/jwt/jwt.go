package jwt

import (
	"time"

	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/shared"
	"github.com/dgrijalva/jwt-go"
)

const (
	keyUserID = "user_id"
)

// Service is a new instance of the JWT service
type Service struct {
	secret string
}

// NewService creates a new instance of the JWT service
func NewService(secret string) Service {
	return Service{
		secret: secret,
	}
}

//GenerateTokenForUserID generates a jwt token and assign a email to it's claims and return it
func (service Service) GenerateTokenForUserID(UserID shared.TransactionID) (result string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)

	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)

	/* Set token claims */
	claims["exp"] = time.Now().AddDate(0, 0, 14)
	claims["nbf"] = time.Now().Unix()
	claims["iat"] = time.Now().Unix()
	claims[keyUserID] = UserID.String()

	result, err = token.SignedString(service.secret)
	if err != nil {
		return result, err
	}

	return result, nil
}

//GetUserIDFromToken parses a jwt token and returns the email it it's claims
func (service Service) GetUserIDFromToken(tokenString string) (userID shared.TransactionID, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return service.secret, nil
	})

	if err != nil {
		return userID, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return userID, err
	}

	return shared.NewTransactionIDFromString(claims[keyUserID].(string))
}
