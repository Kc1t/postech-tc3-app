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
	customerRepo.EXPECT().FindByID(gomock.Any(), "cust-1").Return(customer, nil)
	vehicleRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vehicleRepo := mocks.NewMockVehicleRepository(ctrl)
	customerRepo := mocks.NewMockCustomerRepository(ctrl)
	customerRepo.EXPECT().FindByID(gomock.Any(), "cust-inexistente").Return(nil, domainerrors.ErrNotFound)
	// Create NAO deve ser chamado

	uc := NewCreateVehicle(vehicleRepo, customerRepo)
	vehicle, _ := entities.NewVehicle("cust-inexistente", "ABC-1234", "Fiat", "Uno", 2020)

	err := uc.Execute(context.Background(), vehicle)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("erro = %v, esperava ErrNotFound", err)
	}
}

func TestCreateVehicle_ErroInfraNoFindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	infraErr := errors.New("timeout")

	vehicleRepo := mocks.NewMockVehicleRepository(ctrl)
	customerRepo := mocks.NewMockCustomerRepository(ctrl)
	customerRepo.EXPECT().FindByID(gomock.Any(), "cust-1").Return(nil, infraErr)

	uc := NewCreateVehicle(vehicleRepo, customerRepo)
	vehicle, _ := entities.NewVehicle("cust-1", "ABC-1234", "Fiat", "Uno", 2020)

	err := uc.Execute(context.Background(), vehicle)
	if !errors.Is(err, infraErr) {
		t.Fatalf("erro = %v, esperava %v (erro de infra propagado)", err, infraErr)
	}
}
