package jwt

import (
	"time"

	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/cache"
	"github.com/pkg/errors"

	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/shared"
	"github.com/dgrijalva/jwt-go"
)

const (
	keyUserID = "user_id"
	keyExp    = "exp"
)

var (
	// ErrTokenBlacklisted is thrown when a jwt token is blacklisted
	ErrTokenBlacklisted = errors.New("token has been blacklisted")
)

// Service is a new instance of the JWT service
type Service struct {
	secret string
	cache  cache.Cache
}

// NewService creates a new instance of the JWT service
func NewService(secret string, cache cache.Cache) Service {
	return Service{
		secret: secret,
		cache:  cache,
	}
}

//GenerateTokenForUserID generates a jwt token and assign a email to it's claims and return it
func (service Service) GenerateTokenForUserID(UserID shared.TransactionID) (result string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)

	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)

	/* Set token claims */
	claims[keyExp] = time.Now().AddDate(0, 0, 14)
	claims["nbf"] = time.Now().Unix()
	claims["iat"] = time.Now().Unix()
	claims[keyUserID] = UserID.String()

	result, err = token.SignedString(service.secret)
	if err != nil {
		return result, err
	}

	return result, nil
}

// InvalidateToken invalidates a jwt token
func (service Service) InvalidateToken(tokenString string) (err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return service.secret, nil
	})

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil
	}

	return service.cache.Set(tokenString, "", claims[keyExp].(time.Time).Sub(time.Now()))
}

//GetUserIDFromToken parses a jwt token and returns the email it it's claims
func (service Service) GetUserIDFromToken(tokenString string) (userID shared.TransactionID, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return service.secret, nil
	})

	if err != nil {
		return userID, err
	}

	_, err = service.cache.Get(tokenString)
	if err == nil {
		return userID, ErrTokenBlacklisted
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return userID, err
	}

	return shared.NewTransactionIDFromString(claims[keyUserID].(string))
}
