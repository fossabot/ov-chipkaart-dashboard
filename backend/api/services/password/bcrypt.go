package password

import "golang.org/x/crypto/bcrypt"

// BcryptService is responsible for hashing/comparing passwords using the bcrypt algorithm
type BcryptService struct {
}

// NewBcryptService creates a new instance of the bcrypt service
func NewBcryptService() Service {
	return &BcryptService{}
}

// HashPassword hashes a password
func (service BcryptService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash verifies if the password matches the hashed value
func (service BcryptService) CheckPasswordHash(password string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
