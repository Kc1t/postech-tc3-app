package postgresql

import (
	"context"
	"testing"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/entities"
)

// seedCustomer creates a customer for use as FK in vehicle tests.
func seedCustomer(t *testing.T) *entities.Customer {
	t.Helper()
	repo := NewCustomerRepository(testDB)
	// use document unique per test via t.Name hash or just a fixed doc per file
	c := entities.NewCustomer("Cliente Veiculo", "88888888808", "veiculo@test.com", "11999990009")
	if err := repo.Create(context.Background(), c); err != nil {
		// already exists from a previous run — fetch it
		existing, ferr := repo.FindByDocument(context.Background(), "88888888808")
		if ferr != nil {
			t.Fatalf("seedCustomer failed: %v", err)
		}
		return existing
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Customer{}, "id = ?", c.ID()) })
	return c
}

func TestVehicleRepository_Create(t *testing.T) {
	customer := seedCustomer(t)
	repo := NewVehicleRepository(testDB)
	v := entities.NewVehicle(customer.ID(), "TST0001", "Toyota", "Corolla", 2020)

	if err := repo.Create(context.Background(), v); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if v.ID() == "" {
		t.Fatal("expected ID to be set after create")
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Vehicle{}, "id = ?", v.ID()) })
}

func TestVehicleRepository_Create_DuplicatePlate(t *testing.T) {
	customer := seedCustomer(t)
	repo := NewVehicleRepository(testDB)
	v1 := entities.NewVehicle(customer.ID(), "TST0002", "Honda", "Civic", 2021)
	v2 := entities.NewVehicle(customer.ID(), "TST0002", "Ford", "Ka", 2019)

	if err := repo.Create(context.Background(), v1); err != nil {
		t.Fatalf("first create failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Vehicle{}, "id = ?", v1.ID()) })

	err := repo.Create(context.Background(), v2)
	if err == nil {
		t.Fatal("expected error for duplicate plate")
		t.Cleanup(func() { testDB.Delete(&pgmodel.Vehicle{}, "id = ?", v2.ID()) })
	}
	if !isAlreadyExists(err) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestVehicleRepository_FindByID(t *testing.T) {
	customer := seedCustomer(t)
	repo := NewVehicleRepository(testDB)
	v := entities.NewVehicle(customer.ID(), "TST0003", "Volkswagen", "Golf", 2022)
	if err := repo.Create(context.Background(), v); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Vehicle{}, "id = ?", v.ID()) })

	found, err := repo.FindByID(context.Background(), v.ID())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.Plate() != "TST0003" {
		t.Errorf("expected plate TST0003, got %s", found.Plate())
	}
}

func TestVehicleRepository_FindByID_NotFound(t *testing.T) {
	repo := NewVehicleRepository(testDB)

	_, err := repo.FindByID(context.Background(), "00000000-0000-0000-0000-000000000000")
	if !isDomainNotFound(err) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestVehicleRepository_FindAll(t *testing.T) {
	customer := seedCustomer(t)
	repo := NewVehicleRepository(testDB)
	v := entities.NewVehicle(customer.ID(), "TST0004", "Chevrolet", "Onix", 2023)
	if err := repo.Create(context.Background(), v); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Vehicle{}, "id = ?", v.ID()) })

	all, err := repo.FindAll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(all) == 0 {
		t.Fatal("expected at least one vehicle")
	}
}

func TestVehicleRepository_FindByCustomerID(t *testing.T) {
	customer := seedCustomer(t)
	repo := NewVehicleRepository(testDB)
	v := entities.NewVehicle(customer.ID(), "TST0005", "Fiat", "Uno", 2018)
	if err := repo.Create(context.Background(), v); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Vehicle{}, "id = ?", v.ID()) })

	vehicles, err := repo.FindByCustomerID(context.Background(), customer.ID())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(vehicles) == 0 {
		t.Fatal("expected at least one vehicle for this customer")
	}
}

func TestVehicleRepository_Update(t *testing.T) {
	customer := seedCustomer(t)
	repo := NewVehicleRepository(testDB)
	v := entities.NewVehicle(customer.ID(), "TST0006", "Renault", "Sandero", 2017)
	if err := repo.Create(context.Background(), v); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Vehicle{}, "id = ?", v.ID()) })

	v.SetBrand("Renault Atualizado")
	if err := repo.Update(context.Background(), v); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, _ := repo.FindByID(context.Background(), v.ID())
	if found.Brand() != "Renault Atualizado" {
		t.Errorf("expected updated brand, got %s", found.Brand())
	}
}

func TestVehicleRepository_Delete(t *testing.T) {
	customer := seedCustomer(t)
	repo := NewVehicleRepository(testDB)
	v := entities.NewVehicle(customer.ID(), "TST0007", "Peugeot", "208", 2021)
	if err := repo.Create(context.Background(), v); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if err := repo.Delete(context.Background(), v.ID()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err := repo.FindByID(context.Background(), v.ID())
	if !isDomainNotFound(err) {
		t.Fatal("expected record to be deleted")
	}
}
