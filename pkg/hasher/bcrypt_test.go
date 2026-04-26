package hasher

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestBcryptHasher_Hash_And_Compare(t *testing.T) {
	h := NewBcrypt(bcrypt.MinCost)

	hashed, err := h.Hash("mysecret")
	if err != nil {
		t.Fatalf("Hash: %v", err)
	}
	if hashed == "" {
		t.Fatal("expected non-empty hash")
	}

	if err := h.Compare(hashed, "mysecret"); err != nil {
		t.Fatalf("Compare: expected nil, got %v", err)
	}
}

func TestBcryptHasher_Compare_WrongPassword(t *testing.T) {
	h := NewBcrypt(bcrypt.MinCost)

	hashed, _ := h.Hash("correct")
	if err := h.Compare(hashed, "wrong"); err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}
