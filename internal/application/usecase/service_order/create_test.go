package serviceorderuc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func ft() time.Time { return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) }

func TestCreateServiceOrder_Sucesso(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	vehicle := entities.ReconstituteVehicle("veh-1", "cust-1", "ABC1234", "Fiat", "Uno", 2020, ft(), ft())

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	vehRepo := mocks.NewMockVehicleRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil)
	vehRepo.EXPECT().FindByPlate(gomock.Any(), "ABC1234").Return(vehicle, nil)
	soRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	uc := NewCreateServiceOrder(soRepo, custRepo, vehRepo)

	input := ports.CreateServiceOrderInput{
		CustomerCPF:  "52998224725",
		VehiclePlate: "ABC1234",
	}

	so, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}

	if so.TotalAmount() != 0 {
		t.Errorf("TotalAmount() = %.2f, esperava 0.00", so.TotalAmount())
	}
}

func TestCreateServiceOrder_ErroInfraNoVeiculo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	infraErr := errors.New("timeout")

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	vehRepo := mocks.NewMockVehicleRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil)
	vehRepo.EXPECT().FindByPlate(gomock.Any(), "ABC1234").Return(nil, infraErr)

	uc := NewCreateServiceOrder(soRepo, custRepo, vehRepo)
	input := ports.CreateServiceOrderInput{
		CustomerCPF:  "52998224725",
		VehiclePlate: "ABC1234",
	}

	_, err := uc.Execute(context.Background(), input)
	if !errors.Is(err, infraErr) {
		t.Fatalf("erro = %v, esperava %v (erro de infra propagado)", err, infraErr)
	}
}

func TestCreateServiceOrder_VeiculoNaoEncontrado(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	vehRepo := mocks.NewMockVehicleRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil)
	vehRepo.EXPECT().FindByPlate(gomock.Any(), "XXX0000").Return(nil, domainerrors.ErrNotFound)

	uc := NewCreateServiceOrder(soRepo, custRepo, vehRepo)
	input := ports.CreateServiceOrderInput{
		CustomerCPF:  "52998224725",
		VehiclePlate: "XXX0000",
	}

	_, err := uc.Execute(context.Background(), input)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("erro = %v, esperava ErrNotFound", err)
	}
}

func TestCreateServiceOrder_ErroNaPersistencia(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	vehicle := entities.ReconstituteVehicle("veh-1", "cust-1", "ABC1234", "Fiat", "Uno", 2020, ft(), ft())
	infraErr := errors.New("db write error")

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	vehRepo := mocks.NewMockVehicleRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil)
	vehRepo.EXPECT().FindByPlate(gomock.Any(), "ABC1234").Return(vehicle, nil)
	soRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(infraErr)

	uc := NewCreateServiceOrder(soRepo, custRepo, vehRepo)
	input := ports.CreateServiceOrderInput{
		CustomerCPF:  "52998224725",
		VehiclePlate: "ABC1234",
	}

	_, err := uc.Execute(context.Background(), input)
	if !errors.Is(err, infraErr) {
		t.Fatalf("erro = %v, esperava %v (erro de persistencia propagado)", err, infraErr)
	}
}

func TestCreateServiceOrder_ClienteNaoExiste(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	vehRepo := mocks.NewMockVehicleRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "00000000000").Return(nil, domainerrors.ErrNotFound)

	uc := NewCreateServiceOrder(soRepo, custRepo, vehRepo)
	input := ports.CreateServiceOrderInput{
		CustomerCPF:  "00000000000",
		VehiclePlate: "ABC1234",
	}

	_, err := uc.Execute(context.Background(), input)
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

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil)
	vehRepo.EXPECT().FindByPlate(gomock.Any(), "ABC1234").Return(vehicle, nil)

	uc := NewCreateServiceOrder(soRepo, custRepo, vehRepo)
	input := ports.CreateServiceOrderInput{
		CustomerCPF:  "52998224725",
		VehiclePlate: "ABC1234",
	}

	_, err := uc.Execute(context.Background(), input)
	if !errors.Is(err, domainerrors.ErrVehicleNotFromCustomer) {
		t.Fatalf("erro = %v, esperava ErrVehicleNotFromCustomer", err)
	}
}

func TestCreateServiceOrder_ErroInfraNoCliente(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	infraErr := errors.New("timeout")

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	vehRepo := mocks.NewMockVehicleRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(nil, infraErr)

	uc := NewCreateServiceOrder(soRepo, custRepo, vehRepo)
	input := ports.CreateServiceOrderInput{
		CustomerCPF:  "52998224725",
		VehiclePlate: "ABC1234",
	}

	_, err := uc.Execute(context.Background(), input)
	if !errors.Is(err, infraErr) {
		t.Fatalf("erro = %v, esperava %v (erro de infra propagado)", err, infraErr)
	}
}
