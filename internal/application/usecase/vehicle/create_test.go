package vehicleuc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
)

func fixedTime() time.Time { return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) }

// --- Mocks ---

type mockVehicleRepo struct {
	createFn func(ctx context.Context, v *entities.Vehicle) error
}

func (m *mockVehicleRepo) Create(ctx context.Context, v *entities.Vehicle) error {
	return m.createFn(ctx, v)
}
func (m *mockVehicleRepo) FindByID(ctx context.Context, id string) (*entities.Vehicle, error) {
	return nil, domainerrors.ErrNotFound
}
func (m *mockVehicleRepo) FindByCustomerID(ctx context.Context, customerID string) ([]*entities.Vehicle, error) {
	return nil, nil
}
func (m *mockVehicleRepo) FindAll(ctx context.Context) ([]*entities.Vehicle, error) { return nil, nil }
func (m *mockVehicleRepo) Update(ctx context.Context, v *entities.Vehicle) error    { return nil }
func (m *mockVehicleRepo) Delete(ctx context.Context, id string) error              { return nil }

type mockCustomerRepo struct {
	findByIDFn func(ctx context.Context, id string) (*entities.Customer, error)
}

func (m *mockCustomerRepo) Create(ctx context.Context, c *entities.Customer) error { return nil }
func (m *mockCustomerRepo) FindByID(ctx context.Context, id string) (*entities.Customer, error) {
	return m.findByIDFn(ctx, id)
}
func (m *mockCustomerRepo) FindByDocument(ctx context.Context, doc string) (*entities.Customer, error) {
	return nil, domainerrors.ErrNotFound
}
func (m *mockCustomerRepo) FindAll(ctx context.Context) ([]*entities.Customer, error) {
	return nil, nil
}
func (m *mockCustomerRepo) Update(ctx context.Context, c *entities.Customer) error { return nil }
func (m *mockCustomerRepo) Delete(ctx context.Context, id string) error            { return nil }

// --- Testes ---

func TestCreateVehicle_Sucesso(t *testing.T) {
	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", fixedTime(), fixedTime())

	vehicleRepo := &mockVehicleRepo{
		createFn: func(_ context.Context, _ *entities.Vehicle) error { return nil },
	}
	customerRepo := &mockCustomerRepo{
		findByIDFn: func(_ context.Context, _ string) (*entities.Customer, error) {
			return customer, nil
		},
	}

	uc := NewCreateVehicle(vehicleRepo, customerRepo)
	vehicle, err := entities.NewVehicle("cust-1", "ABC-1234", "Fiat", "Uno", 2020)
	if err != nil {
		t.Fatalf("erro ao criar vehicle: %v", err)
	}

	if err := uc.Execute(context.Background(), vehicle); err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
}

func TestCreateVehicle_ClienteNaoExiste(t *testing.T) {
	vehicleRepo := &mockVehicleRepo{
		createFn: func(_ context.Context, _ *entities.Vehicle) error {
			t.Fatal("Create nao deveria ser chamado")
			return nil
		},
	}
	customerRepo := &mockCustomerRepo{
		findByIDFn: func(_ context.Context, _ string) (*entities.Customer, error) {
			return nil, domainerrors.ErrNotFound
		},
	}

	uc := NewCreateVehicle(vehicleRepo, customerRepo)
	vehicle, _ := entities.NewVehicle("cust-inexistente", "ABC-1234", "Fiat", "Uno", 2020)

	err := uc.Execute(context.Background(), vehicle)
	if err == nil {
		t.Fatal("esperava erro de cliente nao encontrado")
	}
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Errorf("erro = %v, esperava ErrNotFound", err)
	}
}
