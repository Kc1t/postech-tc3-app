package postgresql

import (
	"context"
	"testing"
	"time"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/entities"
)

// seedVehicle creates a vehicle (and the required customer) for service order tests.
func seedVehicle(t *testing.T, plate string) (*entities.Customer, *entities.Vehicle) {
	t.Helper()
	crepo := NewCustomerRepository(testDB)
	doc := "99999999909"
	now := time.Now()
	c := entities.ReconstituteCustomer("", "Cliente OS", doc, "os@test.com", "11999990010", now, now)
	if err := crepo.Create(context.Background(), c); err != nil {
		existing, ferr := crepo.FindByDocument(context.Background(), doc)
		if ferr != nil {
			t.Fatalf("seedVehicle: customer create failed: %v", err)
		}
		c = existing
	} else {
		t.Cleanup(func() { testDB.Delete(&pgmodel.Customer{}, "id = ?", c.ID()) })
	}

	vrepo := NewVehicleRepository(testDB)
	v := entities.ReconstituteVehicle("", c.ID(), plate, "Toyota", "Hilux", 2023, now, now)
	if err := vrepo.Create(context.Background(), v); err != nil {
		t.Fatalf("seedVehicle: vehicle create failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Vehicle{}, "id = ?", v.ID()) })

	return c, v
}

func newTestOrder(customerID, vehicleID string) *entities.ServiceOrder {
	so := entities.NewServiceOrder(customerID, vehicleID)
	return so
}

func TestServiceOrderRepository_Create(t *testing.T) {
	c, v := seedVehicle(t, "OS00001")
	repo := NewServiceOrderRepository(testDB)
	so := newTestOrder(c.ID(), v.ID())

	if err := repo.Create(context.Background(), so); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if so.ID() == "" {
		t.Fatal("expected ID to be set after create")
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.ServiceOrder{}, "id = ?", so.ID()) })
}

func TestServiceOrderRepository_FindByID(t *testing.T) {
	c, v := seedVehicle(t, "OS00002")
	repo := NewServiceOrderRepository(testDB)
	so := newTestOrder(c.ID(), v.ID())
	if err := repo.Create(context.Background(), so); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.ServiceOrder{}, "id = ?", so.ID()) })

	found, err := repo.FindByID(context.Background(), so.ID())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.ID() != so.ID() {
		t.Errorf("expected ID %s, got %s", so.ID(), found.ID())
	}
}

func TestServiceOrderRepository_FindByID_NotFound(t *testing.T) {
	repo := NewServiceOrderRepository(testDB)

	_, err := repo.FindByID(context.Background(), "00000000-0000-0000-0000-000000000000")
	if !isDomainNotFound(err) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestServiceOrderRepository_FindAll(t *testing.T) {
	c, v := seedVehicle(t, "OS00003")
	repo := NewServiceOrderRepository(testDB)
	so := newTestOrder(c.ID(), v.ID())
	if err := repo.Create(context.Background(), so); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.ServiceOrder{}, "id = ?", so.ID()) })

	all, err := repo.FindAll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(all) == 0 {
		t.Fatal("expected at least one service order")
	}
}

func TestServiceOrderRepository_FindByCustomerID(t *testing.T) {
	c, v := seedVehicle(t, "OS00004")
	repo := NewServiceOrderRepository(testDB)
	so := newTestOrder(c.ID(), v.ID())
	if err := repo.Create(context.Background(), so); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.ServiceOrder{}, "id = ?", so.ID()) })

	orders, err := repo.FindByCustomerID(context.Background(), c.ID())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(orders) == 0 {
		t.Fatal("expected at least one service order for this customer")
	}
}

func TestServiceOrderRepository_UpdateStatus(t *testing.T) {
	c, v := seedVehicle(t, "OS00005")
	repo := NewServiceOrderRepository(testDB)
	so := newTestOrder(c.ID(), v.ID())
	if err := repo.Create(context.Background(), so); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.ServiceOrder{}, "id = ?", so.ID()) })

	if err := repo.UpdateStatus(context.Background(), so.ID(), entities.StatusInExecution); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, _ := repo.FindByID(context.Background(), so.ID())
	if found.Status() != entities.StatusInExecution {
		t.Errorf("expected status InProgress, got %s", found.Status())
	}
}

func TestServiceOrderRepository_Update(t *testing.T) {
	c, v := seedVehicle(t, "OS00006")
	repo := NewServiceOrderRepository(testDB)
	so := newTestOrder(c.ID(), v.ID())
	if err := repo.Create(context.Background(), so); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.ServiceOrder{}, "id = ?", so.ID()) })

	so.SetNotes("Observação atualizada")
	if err := repo.Update(context.Background(), so); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, _ := repo.FindByID(context.Background(), so.ID())
	if found.Notes() != "Observação atualizada" {
		t.Errorf("expected updated notes, got %s", found.Notes())
	}
}

func TestServiceOrderRepository_Delete(t *testing.T) {
	c, v := seedVehicle(t, "OS00007")
	repo := NewServiceOrderRepository(testDB)
	so := newTestOrder(c.ID(), v.ID())
	if err := repo.Create(context.Background(), so); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if err := repo.Delete(context.Background(), so.ID()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err := repo.FindByID(context.Background(), so.ID())
	if !isDomainNotFound(err) {
		t.Fatal("expected record to be deleted")
	}
}
