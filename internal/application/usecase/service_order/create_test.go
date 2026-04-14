package serviceorderuc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func ft() time.Time { return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) }

func TestCreateServiceOrder_Sucesso(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	vehicle := entities.ReconstituteVehicle("veh-1", "cust-1", "ABC1234", "Fiat", "Uno", 2020, ft(), ft())
	svc := entities.ReconstituteService("svc-1", "Troca de oleo", "Troca completa", 150.0, 30, ft(), ft())
	part := entities.ReconstitutePart("part-1", "Filtro", "Filtro de oleo", "un", 25.0, 100, ft(), ft())

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	vehRepo := mocks.NewMockVehicleRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	custRepo.EXPECT().FindByID(gomock.Any(), "cust-1").Return(customer, nil)
	vehRepo.EXPECT().FindByID(gomock.Any(), "veh-1").Return(vehicle, nil)
	svcRepo.EXPECT().FindByIDs(gomock.Any(), []string{"svc-1"}).Return([]*entities.Service{svc}, nil)
	partRepo.EXPECT().FindByIDs(gomock.Any(), []string{"part-1"}).Return([]*entities.Part{part}, nil)
	soRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	uc := NewCreateServiceOrder(soRepo, custRepo, vehRepo, svcRepo, partRepo)

	so := entities.NewServiceOrder("cust-1", "veh-1")
	so.AddService(entities.ServiceItem{ServiceID: "svc-1"})
	so.AddPart(entities.PartItem{PartID: "part-1", Quantity: 2})

	if err := uc.Execute(context.Background(), so); err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}

	if so.TotalAmount() != 200.0 { // 150 (servico) + 2*25 (pecas)
		t.Errorf("TotalAmount() = %.2f, esperava 200.00", so.TotalAmount())
	}
	if so.Services()[0].Description != "Troca de oleo" {
		t.Errorf("Service.Description = %q, esperava 'Troca de oleo'", so.Services()[0].Description)
	}
	if so.Services()[0].Price != 150.0 {
		t.Errorf("Service.Price = %.2f, esperava 150.00", so.Services()[0].Price)
	}
}

func TestCreateServiceOrder_SemItens(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	vehicle := entities.ReconstituteVehicle("veh-1", "cust-1", "ABC1234", "Fiat", "Uno", 2020, ft(), ft())

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	vehRepo := mocks.NewMockVehicleRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	custRepo.EXPECT().FindByID(gomock.Any(), "cust-1").Return(customer, nil)
	vehRepo.EXPECT().FindByID(gomock.Any(), "veh-1").Return(vehicle, nil)
	soRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	uc := NewCreateServiceOrder(soRepo, custRepo, vehRepo, svcRepo, partRepo)
	so := entities.NewServiceOrder("cust-1", "veh-1")

	if err := uc.Execute(context.Background(), so); err != nil {
		t.Fatalf("OS sem itens deveria ser permitida: %v", err)
	}
	if so.TotalAmount() != 0 {
		t.Errorf("TotalAmount() = %.2f, esperava 0.00", so.TotalAmount())
	}
}

func TestCreateServiceOrder_ClienteNaoExiste(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	vehRepo := mocks.NewMockVehicleRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	custRepo.EXPECT().FindByID(gomock.Any(), "cust-inexistente").Return(nil, domainerrors.ErrNotFound)

	uc := NewCreateServiceOrder(soRepo, custRepo, vehRepo, svcRepo, partRepo)
	so := entities.NewServiceOrder("cust-inexistente", "veh-1")

	err := uc.Execute(context.Background(), so)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("erro = %v, esperava ErrNotFound", err)
	}
}

func TestCreateServiceOrder_VeiculoNaoPertenceAoCliente(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	vehicle := entities.ReconstituteVehicle("veh-1", "outro-cliente", "ABC1234", "Fiat", "Uno", 2020, ft(), ft())

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	vehRepo := mocks.NewMockVehicleRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	custRepo.EXPECT().FindByID(gomock.Any(), "cust-1").Return(customer, nil)
	vehRepo.EXPECT().FindByID(gomock.Any(), "veh-1").Return(vehicle, nil)

	uc := NewCreateServiceOrder(soRepo, custRepo, vehRepo, svcRepo, partRepo)
	so := entities.NewServiceOrder("cust-1", "veh-1")

	err := uc.Execute(context.Background(), so)
	if !errors.Is(err, domainerrors.ErrVehicleNotFromCustomer) {
		t.Fatalf("erro = %v, esperava ErrVehicleNotFromCustomer", err)
	}
}

func TestCreateServiceOrder_ServicoNaoEncontrado(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	vehicle := entities.ReconstituteVehicle("veh-1", "cust-1", "ABC1234", "Fiat", "Uno", 2020, ft(), ft())

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	vehRepo := mocks.NewMockVehicleRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	custRepo.EXPECT().FindByID(gomock.Any(), "cust-1").Return(customer, nil)
	vehRepo.EXPECT().FindByID(gomock.Any(), "veh-1").Return(vehicle, nil)
	svcRepo.EXPECT().FindByIDs(gomock.Any(), []string{"svc-inexistente"}).Return([]*entities.Service{}, nil)

	uc := NewCreateServiceOrder(soRepo, custRepo, vehRepo, svcRepo, partRepo)
	so := entities.NewServiceOrder("cust-1", "veh-1")
	so.AddService(entities.ServiceItem{ServiceID: "svc-inexistente"})

	err := uc.Execute(context.Background(), so)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("erro = %v, esperava ErrNotFound", err)
	}
}

func TestCreateServiceOrder_EstoqueInsuficiente(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	vehicle := entities.ReconstituteVehicle("veh-1", "cust-1", "ABC1234", "Fiat", "Uno", 2020, ft(), ft())
	partWithLowStock := entities.ReconstitutePart("part-1", "Filtro", "Filtro de oleo", "un", 25.0, 1, ft(), ft())

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	vehRepo := mocks.NewMockVehicleRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	custRepo.EXPECT().FindByID(gomock.Any(), "cust-1").Return(customer, nil)
	vehRepo.EXPECT().FindByID(gomock.Any(), "veh-1").Return(vehicle, nil)
	partRepo.EXPECT().FindByIDs(gomock.Any(), []string{"part-1"}).Return([]*entities.Part{partWithLowStock}, nil)

	uc := NewCreateServiceOrder(soRepo, custRepo, vehRepo, svcRepo, partRepo)
	so := entities.NewServiceOrder("cust-1", "veh-1")
	so.AddPart(entities.PartItem{PartID: "part-1", Quantity: 5}) // pede 5, tem 1

	err := uc.Execute(context.Background(), so)
	if !errors.Is(err, domainerrors.ErrInsufficientStock) {
		t.Fatalf("erro = %v, esperava ErrInsufficientStock", err)
	}
}

func TestCreateServiceOrder_ErroInfraNoCliente(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	infraErr := errors.New("timeout")

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	vehRepo := mocks.NewMockVehicleRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	custRepo.EXPECT().FindByID(gomock.Any(), "cust-1").Return(nil, infraErr)

	uc := NewCreateServiceOrder(soRepo, custRepo, vehRepo, svcRepo, partRepo)
	so := entities.NewServiceOrder("cust-1", "veh-1")

	err := uc.Execute(context.Background(), so)
	if !errors.Is(err, infraErr) {
		t.Fatalf("erro = %v, esperava %v (erro de infra propagado)", err, infraErr)
	}
}
