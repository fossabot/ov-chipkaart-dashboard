package password

// Service is responsible for hashing and comparing hashed passwords
type Service interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password string, hash string) bool
}
