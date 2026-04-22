package vehicleuc

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

func TestCreateVehicle_Sucesso(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now()
	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", now, now)

	vehicleRepo := mocks.NewMockVehicleRepository(ctrl)
	customerRepo := mocks.NewMockCustomerRepository(ctrl)
	customerRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil)
	vehicleRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	uc := NewCreateVehicle(vehicleRepo, customerRepo)
	input := entities.VehicleInput{
		CustomerDocument: "52998224725",
		Plate:            "ABC-1234",
		Brand:            "Fiat",
		Model:            "Uno",
		Year:             2020,
	}

	vehicle, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
	if vehicle == nil || vehicle.CustomerID() != "cust-1" {
		t.Fatalf("vehicle retornado invalido: %+v", vehicle)
	}
}

func TestCreateVehicle_ClienteNaoExiste(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vehicleRepo := mocks.NewMockVehicleRepository(ctrl)
	customerRepo := mocks.NewMockCustomerRepository(ctrl)
	customerRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(nil, domainerrors.ErrNotFound)
	// Create NAO deve ser chamado

	uc := NewCreateVehicle(vehicleRepo, customerRepo)
	input := entities.VehicleInput{
		CustomerDocument: "52998224725",
		Plate:            "ABC-1234",
		Brand:            "Fiat",
		Model:            "Uno",
		Year:             2020,
	}

	_, err := uc.Execute(context.Background(), input)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("erro = %v, esperava ErrNotFound", err)
	}
}

func TestCreateVehicle_ErroInfraNoFindByDocument(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	infraErr := errors.New("timeout")

	vehicleRepo := mocks.NewMockVehicleRepository(ctrl)
	customerRepo := mocks.NewMockCustomerRepository(ctrl)
	customerRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(nil, infraErr)

	uc := NewCreateVehicle(vehicleRepo, customerRepo)
	input := entities.VehicleInput{
		CustomerDocument: "52998224725",
		Plate:            "ABC-1234",
		Brand:            "Fiat",
		Model:            "Uno",
		Year:             2020,
	}

	_, err := uc.Execute(context.Background(), input)
	if !errors.Is(err, infraErr) {
		t.Fatalf("erro = %v, esperava %v (erro de infra propagado)", err, infraErr)
	}
}

func TestCreateVehicle_PlacaInvalida(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now()
	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", now, now)

	vehicleRepo := mocks.NewMockVehicleRepository(ctrl)
	customerRepo := mocks.NewMockCustomerRepository(ctrl)
	customerRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil)
	// Create NAO deve ser chamado

	uc := NewCreateVehicle(vehicleRepo, customerRepo)
	input := entities.VehicleInput{
		CustomerDocument: "52998224725",
		Plate:            "INVALID",
		Brand:            "Fiat",
		Model:            "Uno",
		Year:             2020,
	}

	_, err := uc.Execute(context.Background(), input)
	if !errors.Is(err, domainerrors.ErrInvalidPlate) {
		t.Fatalf("erro = %v, esperava ErrInvalidPlate", err)
	}
}
