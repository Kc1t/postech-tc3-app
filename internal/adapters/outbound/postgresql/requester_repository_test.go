package postgresql

import (
	"context"
	"testing"
	"time"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
)

func TestRequesterRepository_Create(t *testing.T) {
	repo := NewRequesterRepository(testDB)
	now := time.Now()
	c := entities.ReconstituteRequester("", "João Silva", "11111111101", "joao@test.com", "11999990001", now, now)

	err := repo.Create(context.Background(), c)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if c.ID() == "" {
		t.Fatal("expected ID to be set after create")
	}

	t.Cleanup(func() {
		testDB.Delete(&pgmodel.Requester{}, "id = ?", c.ID())
	})
}

func TestRequesterRepository_Create_DuplicateDocument(t *testing.T) {
	repo := NewRequesterRepository(testDB)
	now := time.Now()
	c1 := entities.ReconstituteRequester("", "João", "22222222202", "j1@test.com", "11999990002", now, now)
	c2 := entities.ReconstituteRequester("", "João2", "22222222202", "j2@test.com", "11999990003", now, now)

	if err := repo.Create(context.Background(), c1); err != nil {
		t.Fatalf("first create failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Requester{}, "id = ?", c1.ID()) })

	err := repo.Create(context.Background(), c2)
	if err == nil {
		t.Cleanup(func() { testDB.Delete(&pgmodel.Requester{}, "id = ?", c2.ID()) })
		t.Fatal("expected error for duplicate document")
	}
	if !isAlreadyExists(err) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestRequesterRepository_FindByID(t *testing.T) {
	repo := NewRequesterRepository(testDB)
	now := time.Now()
	c := entities.ReconstituteRequester("", "Maria", "33333333303", "maria@test.com", "11999990004", now, now)
	if err := repo.Create(context.Background(), c); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Requester{}, "id = ?", c.ID()) })

	found, err := repo.FindByID(context.Background(), c.ID())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.Name() != "Maria" {
		t.Errorf("expected name Maria, got %s", found.Name())
	}
}

func TestRequesterRepository_FindByID_NotFound(t *testing.T) {
	repo := NewRequesterRepository(testDB)

	_, err := repo.FindByID(context.Background(), "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Fatal("expected ErrNotFound")
	}
	if !isDomainNotFound(err) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRequesterRepository_FindByDocument(t *testing.T) {
	repo := NewRequesterRepository(testDB)
	now := time.Now()
	c := entities.ReconstituteRequester("", "Carlos", "44444444404", "carlos@test.com", "11999990005", now, now)
	if err := repo.Create(context.Background(), c); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Requester{}, "id = ?", c.ID()) })

	found, err := repo.FindByDocument(context.Background(), "44444444404")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.Document() != "44444444404" {
		t.Errorf("expected document 44444444404, got %s", found.Document())
	}
}

func TestRequesterRepository_FindByDocument_NotFound(t *testing.T) {
	repo := NewRequesterRepository(testDB)

	_, err := repo.FindByDocument(context.Background(), "00000000000")
	if err == nil {
		t.Fatal("expected ErrNotFound")
	}
	if !isDomainNotFound(err) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRequesterRepository_FindAll(t *testing.T) {
	repo := NewRequesterRepository(testDB)
	now := time.Now()
	c := entities.ReconstituteRequester("", "Ana", "55555555505", "ana@test.com", "11999990006", now, now)
	if err := repo.Create(context.Background(), c); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Requester{}, "id = ?", c.ID()) })

	all, err := repo.FindAll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(all) == 0 {
		t.Fatal("expected at least one requester")
	}
}

func TestRequesterRepository_Update(t *testing.T) {
	repo := NewRequesterRepository(testDB)
	now := time.Now()
	c := entities.ReconstituteRequester("", "Pedro", "66666666606", "pedro@test.com", "11999990007", now, now)
	if err := repo.Create(context.Background(), c); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Requester{}, "id = ?", c.ID()) })

	c.SetName("Pedro Atualizado")
	if err := repo.Update(context.Background(), c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, _ := repo.FindByID(context.Background(), c.ID())
	if found.Name() != "Pedro Atualizado" {
		t.Errorf("expected updated name, got %s", found.Name())
	}
}

func TestRequesterRepository_Delete(t *testing.T) {
	repo := NewRequesterRepository(testDB)
	now := time.Now()
	c := entities.ReconstituteRequester("", "Lucas", "77777777707", "lucas@test.com", "11999990008", now, now)
	if err := repo.Create(context.Background(), c); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if err := repo.Delete(context.Background(), c.ID()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err := repo.FindByID(context.Background(), c.ID())
	if !isDomainNotFound(err) {
		t.Fatal("expected record to be deleted")
	}
}

func isAlreadyExists(err error) bool {
	return err != nil && err.Error() == domainerrors.ErrAlreadyExists.Error()
}

func isDomainNotFound(err error) bool {
	return err != nil && err.Error() == domainerrors.ErrNotFound.Error()
}
