package vehicleuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestCreateVehicle_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.NewCustomer("João", "12345678901", "j@j.com", "11999")
	customer.SetID("cust-1")

	vehicleRepo := mocks.NewMockVehicleRepository(ctrl)
	customerRepo := mocks.NewMockCustomerRepository(ctrl)
	customerRepo.EXPECT().FindByDocument(gomock.Any(), "12345678901").Return(customer, nil)
	vehicleRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	v := entities.NewVehicle("", "ABC1234", "Toyota", "Corolla", 2020)
	uc := NewCreateVehicle(vehicleRepo, customerRepo)
	if err := uc.Execute(context.Background(), "12345678901", v); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if v.CustomerID() != "cust-1" {
		t.Errorf("expected customerID %q, got %q", "cust-1", v.CustomerID())
	}
}

func TestCreateVehicle_Execute_CustomerNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vehicleRepo := mocks.NewMockVehicleRepository(ctrl)
	customerRepo := mocks.NewMockCustomerRepository(ctrl)
	customerRepo.EXPECT().FindByDocument(gomock.Any(), "00000000000").Return(nil, domainerrors.ErrNotFound)

	v := entities.NewVehicle("", "ABC1234", "Toyota", "Corolla", 2020)
	uc := NewCreateVehicle(vehicleRepo, customerRepo)
	err := uc.Execute(context.Background(), "00000000000", v)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCreateVehicle_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.NewCustomer("João", "12345678901", "j@j.com", "11999")
	customer.SetID("cust-1")
	repoErr := errors.New("db error")

	vehicleRepo := mocks.NewMockVehicleRepository(ctrl)
	customerRepo := mocks.NewMockCustomerRepository(ctrl)
	customerRepo.EXPECT().FindByDocument(gomock.Any(), "12345678901").Return(customer, nil)
	vehicleRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(repoErr)

	v := entities.NewVehicle("", "ABC1234", "Toyota", "Corolla", 2020)
	uc := NewCreateVehicle(vehicleRepo, customerRepo)
	if err := uc.Execute(context.Background(), "12345678901", v); !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
