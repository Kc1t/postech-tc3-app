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

func TestGetVehicle_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := entities.NewVehicle("cust-1", "ABC1234", "Toyota", "Corolla", 2020)
	repo := mocks.NewMockVehicleRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "veh-1").Return(expected, nil)

	uc := NewGetVehicle(repo)
	got, err := uc.Execute(context.Background(), "veh-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got != expected {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestGetVehicle_Execute_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockVehicleRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "veh-x").Return(nil, domainerrors.ErrNotFound)

	uc := NewGetVehicle(repo)
	_, err := uc.Execute(context.Background(), "veh-x")
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
