package postgresql

import (
	"context"
	"errors"
	"testing"
	"time"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
)

// newTestUser cria um usuario direto no DB pra servir como dono dos refresh tokens
// nos testes. Retorna ID e uma funcao de cleanup.
func newTestUser(t *testing.T, email string) (string, func()) {
	t.Helper()
	u := entities.NewUser("Test User", email, "hashed", entities.RoleClient)
	repo := NewUserRepository(testDB)
	if err := repo.Create(context.Background(), u); err != nil {
		t.Fatalf("setup user failed: %v", err)
	}
	cleanup := func() { testDB.Delete(&pgmodel.User{}, "id = ?", u.ID()) }
	return u.ID(), cleanup
}

func TestRefreshTokenRepository_Create(t *testing.T) {
	userID, cleanupUser := newTestUser(t, "rt-create@test.com")
	t.Cleanup(cleanupUser)

	repo := NewRefreshTokenRepository(testDB)
	rt := entities.NewRefreshToken(userID, "hash-create-1", time.Now().Add(24*time.Hour))

	if err := repo.Create(context.Background(), rt); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rt.ID() == "" {
		t.Fatal("expected ID to be set by DB after create")
	}

	t.Cleanup(func() { testDB.Delete(&pgmodel.RefreshToken{}, "id = ?", rt.ID()) })
}

func TestRefreshTokenRepository_FindByTokenHash(t *testing.T) {
	userID, cleanupUser := newTestUser(t, "rt-find@test.com")
	t.Cleanup(cleanupUser)

	repo := NewRefreshTokenRepository(testDB)
	rt := entities.NewRefreshToken(userID, "hash-find-1", time.Now().Add(24*time.Hour))
	if err := repo.Create(context.Background(), rt); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.RefreshToken{}, "id = ?", rt.ID()) })

	found, err := repo.FindByTokenHash(context.Background(), "hash-find-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.ID() != rt.ID() {
		t.Errorf("expected ID %s, got %s", rt.ID(), found.ID())
	}
	if found.UserID() != userID {
		t.Errorf("expected userID %s, got %s", userID, found.UserID())
	}
}

func TestRefreshTokenRepository_FindByTokenHash_NotFound(t *testing.T) {
	repo := NewRefreshTokenRepository(testDB)

	_, err := repo.FindByTokenHash(context.Background(), "nao-existe")
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRefreshTokenRepository_FindByTokenHash_IgnoresRevoked(t *testing.T) {
	userID, cleanupUser := newTestUser(t, "rt-revoked@test.com")
	t.Cleanup(cleanupUser)

	repo := NewRefreshTokenRepository(testDB)
	rt := entities.NewRefreshToken(userID, "hash-revoked-1", time.Now().Add(24*time.Hour))
	if err := repo.Create(context.Background(), rt); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.RefreshToken{}, "id = ?", rt.ID()) })

	if err := repo.Revoke(context.Background(), rt.ID()); err != nil {
		t.Fatalf("revoke failed: %v", err)
	}

	_, err := repo.FindByTokenHash(context.Background(), "hash-revoked-1")
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for revoked token, got %v", err)
	}
}

func TestRefreshTokenRepository_Revoke(t *testing.T) {
	userID, cleanupUser := newTestUser(t, "rt-revoke@test.com")
	t.Cleanup(cleanupUser)

	repo := NewRefreshTokenRepository(testDB)
	rt := entities.NewRefreshToken(userID, "hash-revoke-1", time.Now().Add(24*time.Hour))
	if err := repo.Create(context.Background(), rt); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.RefreshToken{}, "id = ?", rt.ID()) })

	if err := repo.Revoke(context.Background(), rt.ID()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Segunda chamada precisa falhar: token ja revogado nao pode ser revogado de novo.
	err := repo.Revoke(context.Background(), rt.ID())
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("expected ErrNotFound on second revoke (single-use), got %v", err)
	}
}

func TestRefreshTokenRepository_Revoke_NotFound(t *testing.T) {
	repo := NewRefreshTokenRepository(testDB)

	err := repo.Revoke(context.Background(), "00000000-0000-0000-0000-000000000000")
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRefreshTokenRepository_RevokeByUserID(t *testing.T) {
	userID, cleanupUser := newTestUser(t, "rt-revoke-all@test.com")
	t.Cleanup(cleanupUser)

	repo := NewRefreshTokenRepository(testDB)
	rt1 := entities.NewRefreshToken(userID, "hash-all-1", time.Now().Add(24*time.Hour))
	rt2 := entities.NewRefreshToken(userID, "hash-all-2", time.Now().Add(24*time.Hour))
	if err := repo.Create(context.Background(), rt1); err != nil {
		t.Fatalf("setup rt1 failed: %v", err)
	}
	if err := repo.Create(context.Background(), rt2); err != nil {
		t.Fatalf("setup rt2 failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.RefreshToken{}, "user_id = ?", userID) })

	if err := repo.RevokeByUserID(context.Background(), userID); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Nenhum dos dois tokens pode mais ser encontrado (filtro ignora revoked = true).
	if _, err := repo.FindByTokenHash(context.Background(), "hash-all-1"); !errors.Is(err, domainerrors.ErrNotFound) {
		t.Errorf("expected rt1 to be unreachable, got %v", err)
	}
	if _, err := repo.FindByTokenHash(context.Background(), "hash-all-2"); !errors.Is(err, domainerrors.ErrNotFound) {
		t.Errorf("expected rt2 to be unreachable, got %v", err)
	}
}

func TestRefreshTokenRepository_RotateToken(t *testing.T) {
	userID, cleanupUser := newTestUser(t, "rt-rotate@test.com")
	t.Cleanup(cleanupUser)

	repo := NewRefreshTokenRepository(testDB)
	old := entities.NewRefreshToken(userID, "hash-rotate-old", time.Now().Add(24*time.Hour))
	if err := repo.Create(context.Background(), old); err != nil {
		t.Fatalf("setup old failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.RefreshToken{}, "user_id = ?", userID) })

	newRT := entities.NewRefreshToken(userID, "hash-rotate-new", time.Now().Add(24*time.Hour))
	if err := repo.RotateToken(context.Background(), old.ID(), newRT); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if newRT.ID() == "" {
		t.Fatal("expected new token ID to be set by DB after rotate")
	}

	// Old ficou revogado (sumiu do find).
	if _, err := repo.FindByTokenHash(context.Background(), "hash-rotate-old"); !errors.Is(err, domainerrors.ErrNotFound) {
		t.Errorf("expected old to be revoked, got %v", err)
	}
	// New continua valido.
	found, err := repo.FindByTokenHash(context.Background(), "hash-rotate-new")
	if err != nil {
		t.Fatalf("expected new to be found, got %v", err)
	}
	if found.ID() != newRT.ID() {
		t.Errorf("expected new ID %s, got %s", newRT.ID(), found.ID())
	}
}

// TestRefreshTokenRepository_RotateToken_SingleUse garante que a rotacao e atomica:
// uma segunda tentativa de rotar o mesmo old token falha. Simula o cenario de reuse
// attack — replay do mesmo refresh token duas vezes.
func TestRefreshTokenRepository_RotateToken_SingleUse(t *testing.T) {
	userID, cleanupUser := newTestUser(t, "rt-single-use@test.com")
	t.Cleanup(cleanupUser)

	repo := NewRefreshTokenRepository(testDB)
	old := entities.NewRefreshToken(userID, "hash-single-old", time.Now().Add(24*time.Hour))
	if err := repo.Create(context.Background(), old); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.RefreshToken{}, "user_id = ?", userID) })

	first := entities.NewRefreshToken(userID, "hash-single-first", time.Now().Add(24*time.Hour))
	if err := repo.RotateToken(context.Background(), old.ID(), first); err != nil {
		t.Fatalf("first rotate failed: %v", err)
	}

	// Replay: usar o mesmo old.ID() uma segunda vez deve falhar com ErrNotFound
	// (o WHERE revoked = false nao encontra o registro).
	second := entities.NewRefreshToken(userID, "hash-single-second", time.Now().Add(24*time.Hour))
	err := repo.RotateToken(context.Background(), old.ID(), second)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("expected ErrNotFound on second rotate (single-use), got %v", err)
	}
}

func TestRefreshTokenRepository_RotateToken_NotFound(t *testing.T) {
	userID, cleanupUser := newTestUser(t, "rt-rotate-nf@test.com")
	t.Cleanup(cleanupUser)

	repo := NewRefreshTokenRepository(testDB)
	newRT := entities.NewRefreshToken(userID, "hash-rotate-nf", time.Now().Add(24*time.Hour))

	err := repo.RotateToken(context.Background(), "00000000-0000-0000-0000-000000000000", newRT)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
