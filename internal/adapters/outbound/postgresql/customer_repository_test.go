package postgresql

import (
	"context"
	"testing"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/domain/entities"
)

func TestCustomerRepository_Create(t *testing.T) {
	repo := NewCustomerRepository(testDB)
	c := entities.NewCustomer("João Silva", "11111111101", "joao@test.com", "11999990001")

	err := repo.Create(context.Background(), c)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if c.ID() == "" {
		t.Fatal("expected ID to be set after create")
	}

	t.Cleanup(func() {
		testDB.Delete(&pgmodel.Customer{}, "id = ?", c.ID())
	})
}

func TestCustomerRepository_Create_DuplicateDocument(t *testing.T) {
	repo := NewCustomerRepository(testDB)
	c1 := entities.NewCustomer("João", "22222222202", "j1@test.com", "11999990002")
	c2 := entities.NewCustomer("João2", "22222222202", "j2@test.com", "11999990003")

	if err := repo.Create(context.Background(), c1); err != nil {
		t.Fatalf("first create failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Customer{}, "id = ?", c1.ID()) })

	err := repo.Create(context.Background(), c2)
	if err == nil {
		t.Fatal("expected error for duplicate document")
		t.Cleanup(func() { testDB.Delete(&pgmodel.Customer{}, "id = ?", c2.ID()) })
	}
	if !isAlreadyExists(err) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestCustomerRepository_FindByID(t *testing.T) {
	repo := NewCustomerRepository(testDB)
	c := entities.NewCustomer("Maria", "33333333303", "maria@test.com", "11999990004")
	if err := repo.Create(context.Background(), c); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Customer{}, "id = ?", c.ID()) })

	found, err := repo.FindByID(context.Background(), c.ID())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.Name() != "Maria" {
		t.Errorf("expected name Maria, got %s", found.Name())
	}
}

func TestCustomerRepository_FindByID_NotFound(t *testing.T) {
	repo := NewCustomerRepository(testDB)

	_, err := repo.FindByID(context.Background(), "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Fatal("expected ErrNotFound")
	}
	if !isDomainNotFound(err) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCustomerRepository_FindByDocument(t *testing.T) {
	repo := NewCustomerRepository(testDB)
	c := entities.NewCustomer("Carlos", "44444444404", "carlos@test.com", "11999990005")
	if err := repo.Create(context.Background(), c); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Customer{}, "id = ?", c.ID()) })

	found, err := repo.FindByDocument(context.Background(), "44444444404")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.Document() != "44444444404" {
		t.Errorf("expected document 44444444404, got %s", found.Document())
	}
}

func TestCustomerRepository_FindByDocument_NotFound(t *testing.T) {
	repo := NewCustomerRepository(testDB)

	_, err := repo.FindByDocument(context.Background(), "00000000000")
	if err == nil {
		t.Fatal("expected ErrNotFound")
	}
	if !isDomainNotFound(err) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCustomerRepository_FindAll(t *testing.T) {
	repo := NewCustomerRepository(testDB)
	c := entities.NewCustomer("Ana", "55555555505", "ana@test.com", "11999990006")
	if err := repo.Create(context.Background(), c); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Customer{}, "id = ?", c.ID()) })

	all, err := repo.FindAll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(all) == 0 {
		t.Fatal("expected at least one customer")
	}
}

func TestCustomerRepository_Update(t *testing.T) {
	repo := NewCustomerRepository(testDB)
	c := entities.NewCustomer("Pedro", "66666666606", "pedro@test.com", "11999990007")
	if err := repo.Create(context.Background(), c); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Customer{}, "id = ?", c.ID()) })

	c.SetName("Pedro Atualizado")
	if err := repo.Update(context.Background(), c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, _ := repo.FindByID(context.Background(), c.ID())
	if found.Name() != "Pedro Atualizado" {
		t.Errorf("expected updated name, got %s", found.Name())
	}
}

func TestCustomerRepository_Delete(t *testing.T) {
	repo := NewCustomerRepository(testDB)
	c := entities.NewCustomer("Lucas", "77777777707", "lucas@test.com", "11999990008")
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

// helpers

func isAlreadyExists(err error) bool {
	return err != nil && err.Error() == domainerrors.ErrAlreadyExists.Error()
}

func isDomainNotFound(err error) bool {
	return err != nil && err.Error() == domainerrors.ErrNotFound.Error()
}
