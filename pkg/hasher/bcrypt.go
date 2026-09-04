// Package hasher fornece implementacao de hashing de senhas via bcrypt.
// A interface PasswordHasher fica em ports/, aqui apenas a implementacao.
package hasher

import "golang.org/x/crypto/bcrypt"

type bcryptHasher struct {
	cost int
}

// NewBcrypt cria um PasswordHasher usando bcrypt com o cost especificado.
func NewBcrypt(cost int) *bcryptHasher {
	return &bcryptHasher{cost: cost}
}

func (h *bcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (h *bcryptHasher) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
