package serviceorderuc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
)

// --- Mocks ---

type mockServiceOrderRepo struct {
	createFn func(ctx context.Context, so *entities.ServiceOrder) error
}

func (m *mockServiceOrderRepo) Create(ctx context.Context, so *entities.ServiceOrder) error {
	return m.createFn(ctx, so)
}
func (m *mockServiceOrderRepo) FindByID(ctx context.Context, id string) (*entities.ServiceOrder, error) {
	return nil, domainerrors.ErrNotFound
}
func (m *mockServiceOrderRepo) FindAll(ctx context.Context) ([]*entities.ServiceOrder, error) {
	return nil, nil
}
func (m *mockServiceOrderRepo) FindByCustomerID(ctx context.Context, id string) ([]*entities.ServiceOrder, error) {
	return nil, nil
}
func (m *mockServiceOrderRepo) UpdateStatus(ctx context.Context, id string, s entities.OrderStatus) error {
	return nil
}
func (m *mockServiceOrderRepo) Update(ctx context.Context, so *entities.ServiceOrder) error {
	return nil
}
func (m *mockServiceOrderRepo) Delete(ctx context.Context, id string) error { return nil }

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

type mockVehicleRepo struct {
	findByIDFn func(ctx context.Context, id string) (*entities.Vehicle, error)
}

func (m *mockVehicleRepo) Create(ctx context.Context, v *entities.Vehicle) error { return nil }
func (m *mockVehicleRepo) FindByID(ctx context.Context, id string) (*entities.Vehicle, error) {
	return m.findByIDFn(ctx, id)
}
func (m *mockVehicleRepo) FindByCustomerID(ctx context.Context, id string) ([]*entities.Vehicle, error) {
	return nil, nil
}
func (m *mockVehicleRepo) FindAll(ctx context.Context) ([]*entities.Vehicle, error) { return nil, nil }
func (m *mockVehicleRepo) Update(ctx context.Context, v *entities.Vehicle) error    { return nil }
func (m *mockVehicleRepo) Delete(ctx context.Context, id string) error              { return nil }

type mockServiceRepo struct {
	findByIDsFn func(ctx context.Context, ids []string) ([]*entities.Service, error)
}

func (m *mockServiceRepo) Create(ctx context.Context, s *entities.Service) error { return nil }
func (m *mockServiceRepo) FindByID(ctx context.Context, id string) (*entities.Service, error) {
	return nil, domainerrors.ErrNotFound
}
func (m *mockServiceRepo) FindByIDs(ctx context.Context, ids []string) ([]*entities.Service, error) {
	return m.findByIDsFn(ctx, ids)
}
func (m *mockServiceRepo) FindAll(ctx context.Context) ([]*entities.Service, error) { return nil, nil }
func (m *mockServiceRepo) Update(ctx context.Context, s *entities.Service) error    { return nil }
func (m *mockServiceRepo) Delete(ctx context.Context, id string) error              { return nil }

type mockPartRepo struct {
	findByIDsFn func(ctx context.Context, ids []string) ([]*entities.Part, error)
}

func (m *mockPartRepo) Create(ctx context.Context, p *entities.Part) error { return nil }
func (m *mockPartRepo) FindByID(ctx context.Context, id string) (*entities.Part, error) {
	return nil, domainerrors.ErrNotFound
}
func (m *mockPartRepo) FindByIDs(ctx context.Context, ids []string) ([]*entities.Part, error) {
	return m.findByIDsFn(ctx, ids)
}
func (m *mockPartRepo) FindAll(ctx context.Context) ([]*entities.Part, error) { return nil, nil }
func (m *mockPartRepo) Update(ctx context.Context, p *entities.Part) error    { return nil }
func (m *mockPartRepo) Delete(ctx context.Context, id string) error           { return nil }
func (m *mockPartRepo) UpdateStock(ctx context.Context, id string, delta int) error {
	return nil
}

// --- Helpers ---

func ft() time.Time { return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) }

func defaultCustomerRepo() *mockCustomerRepo {
	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	return &mockCustomerRepo{
		findByIDFn: func(_ context.Context, _ string) (*entities.Customer, error) {
			return customer, nil
		},
	}
}

func defaultVehicleRepo(customerID string) *mockVehicleRepo {
	vehicle := entities.ReconstituteVehicle("veh-1", customerID, "ABC1234", "Fiat", "Uno", 2020, ft(), ft())
	return &mockVehicleRepo{
		findByIDFn: func(_ context.Context, _ string) (*entities.Vehicle, error) {
			return vehicle, nil
		},
	}
}

func defaultServiceRepo() *mockServiceRepo {
	svc := entities.ReconstituteService("svc-1", "Troca de oleo", "Troca completa", 150.0, 30, ft(), ft())
	return &mockServiceRepo{
		findByIDsFn: func(_ context.Context, ids []string) ([]*entities.Service, error) {
			result := make([]*entities.Service, 0)
			for range ids {
				result = append(result, svc)
			}
			return result, nil
		},
	}
}

func defaultPartRepo() *mockPartRepo {
	part := entities.ReconstitutePart("part-1", "Filtro", "Filtro de oleo", "un", 25.0, 100, ft(), ft())
	return &mockPartRepo{
		findByIDsFn: func(_ context.Context, ids []string) ([]*entities.Part, error) {
			result := make([]*entities.Part, 0)
			for range ids {
				result = append(result, part)
			}
			return result, nil
		},
	}
}

// --- Testes ---

func TestCreateServiceOrder_Sucesso(t *testing.T) {
	soRepo := &mockServiceOrderRepo{
		createFn: func(_ context.Context, _ *entities.ServiceOrder) error { return nil },
	}

	uc := NewCreateServiceOrder(soRepo, defaultCustomerRepo(), defaultVehicleRepo("cust-1"), defaultServiceRepo(), defaultPartRepo())

	so := entities.NewServiceOrder("cust-1", "veh-1")
	so.AddService(entities.ServiceItem{ServiceID: "svc-1"})
	so.AddPart(entities.PartItem{PartID: "part-1", Quantity: 2})

	if err := uc.Execute(context.Background(), so); err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}

	// Verificar que precos foram preenchidos pelo orcamento automatico
	if so.TotalAmount() != 200.0 { // 150 (servico) + 2*25 (pecas)
		t.Errorf("TotalAmount() = %.2f, esperava 200.00", so.TotalAmount())
	}

	// Verificar que descricao veio do cadastro
	if so.Services()[0].Description != "Troca de oleo" {
		t.Errorf("Service.Description = %q, esperava 'Troca de oleo'", so.Services()[0].Description)
	}
	if so.Services()[0].Price != 150.0 {
		t.Errorf("Service.Price = %.2f, esperava 150.00", so.Services()[0].Price)
	}
}

func TestCreateServiceOrder_SemItens(t *testing.T) {
	soRepo := &mockServiceOrderRepo{
		createFn: func(_ context.Context, _ *entities.ServiceOrder) error { return nil },
	}

	uc := NewCreateServiceOrder(soRepo, defaultCustomerRepo(), defaultVehicleRepo("cust-1"), defaultServiceRepo(), defaultPartRepo())
	so := entities.NewServiceOrder("cust-1", "veh-1")

	if err := uc.Execute(context.Background(), so); err != nil {
		t.Fatalf("OS sem itens deveria ser permitida: %v", err)
	}
	if so.TotalAmount() != 0 {
		t.Errorf("TotalAmount() = %.2f, esperava 0.00", so.TotalAmount())
	}
}

func TestCreateServiceOrder_ClienteNaoExiste(t *testing.T) {
	soRepo := &mockServiceOrderRepo{
		createFn: func(_ context.Context, _ *entities.ServiceOrder) error {
			t.Fatal("Create nao deveria ser chamado")
			return nil
		},
	}
	custRepo := &mockCustomerRepo{
		findByIDFn: func(_ context.Context, _ string) (*entities.Customer, error) {
			return nil, domainerrors.ErrNotFound
		},
	}

	uc := NewCreateServiceOrder(soRepo, custRepo, defaultVehicleRepo("cust-1"), defaultServiceRepo(), defaultPartRepo())
	so := entities.NewServiceOrder("cust-inexistente", "veh-1")

	err := uc.Execute(context.Background(), so)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Errorf("erro = %v, esperava ErrNotFound", err)
	}
}

func TestCreateServiceOrder_VeiculoNaoPertenceAoCliente(t *testing.T) {
	soRepo := &mockServiceOrderRepo{
		createFn: func(_ context.Context, _ *entities.ServiceOrder) error {
			t.Fatal("Create nao deveria ser chamado")
			return nil
		},
	}
	// Veiculo pertence a outro cliente
	vehRepo := defaultVehicleRepo("outro-cliente")

	uc := NewCreateServiceOrder(soRepo, defaultCustomerRepo(), vehRepo, defaultServiceRepo(), defaultPartRepo())
	so := entities.NewServiceOrder("cust-1", "veh-1")

	err := uc.Execute(context.Background(), so)
	if !errors.Is(err, domainerrors.ErrVehicleNotFromCustomer) {
		t.Errorf("erro = %v, esperava ErrVehicleNotFromCustomer", err)
	}
}

func TestCreateServiceOrder_ServicoNaoEncontrado(t *testing.T) {
	soRepo := &mockServiceOrderRepo{
		createFn: func(_ context.Context, _ *entities.ServiceOrder) error { return nil },
	}
	svcRepo := &mockServiceRepo{
		findByIDsFn: func(_ context.Context, _ []string) ([]*entities.Service, error) {
			return []*entities.Service{}, nil // retorna vazio = servico nao encontrado
		},
	}

	uc := NewCreateServiceOrder(soRepo, defaultCustomerRepo(), defaultVehicleRepo("cust-1"), svcRepo, defaultPartRepo())
	so := entities.NewServiceOrder("cust-1", "veh-1")
	so.AddService(entities.ServiceItem{ServiceID: "svc-inexistente"})

	err := uc.Execute(context.Background(), so)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Errorf("erro = %v, esperava ErrNotFound", err)
	}
}

func TestCreateServiceOrder_EstoqueInsuficiente(t *testing.T) {
	soRepo := &mockServiceOrderRepo{
		createFn: func(_ context.Context, _ *entities.ServiceOrder) error { return nil },
	}
	// Peca com estoque = 1
	partWithLowStock := entities.ReconstitutePart("part-1", "Filtro", "Filtro de oleo", "un", 25.0, 1, ft(), ft())
	partRepo := &mockPartRepo{
		findByIDsFn: func(_ context.Context, ids []string) ([]*entities.Part, error) {
			return []*entities.Part{partWithLowStock}, nil
		},
	}

	uc := NewCreateServiceOrder(soRepo, defaultCustomerRepo(), defaultVehicleRepo("cust-1"), defaultServiceRepo(), partRepo)
	so := entities.NewServiceOrder("cust-1", "veh-1")
	so.AddPart(entities.PartItem{PartID: "part-1", Quantity: 5}) // pede 5, tem 1

	err := uc.Execute(context.Background(), so)
	if !errors.Is(err, domainerrors.ErrInsufficientStock) {
		t.Errorf("erro = %v, esperava ErrInsufficientStock", err)
	}
}
