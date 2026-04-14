package vehicleuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestListVehiclesByCustomer_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := []*entities.Vehicle{
		entities.NewVehicle("cust-1", "ABC1234", "Toyota", "Corolla", 2020),
	}
	repo := mocks.NewMockVehicleRepository(ctrl)
	repo.EXPECT().FindByCustomerID(gomock.Any(), "cust-1").Return(expected, nil)

	uc := NewListVehiclesByCustomer(repo)
	got, err := uc.Execute(context.Background(), "cust-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 vehicle, got %d", len(got))
	}
}

func TestListVehiclesByCustomer_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	repo := mocks.NewMockVehicleRepository(ctrl)
	repo.EXPECT().FindByCustomerID(gomock.Any(), "cust-1").Return(nil, repoErr)

	uc := NewListVehiclesByCustomer(repo)
	_, err := uc.Execute(context.Background(), "cust-1")
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
