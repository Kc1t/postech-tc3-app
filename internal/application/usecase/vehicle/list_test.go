package vehicleuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestListVehicles_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := []*entities.Vehicle{
		entities.NewVehicle("cust-1", "ABC1234", "Toyota", "Corolla", 2020),
	}
	repo := mocks.NewMockVehicleRepository(ctrl)
	repo.EXPECT().FindAll(gomock.Any()).Return(expected, nil)

	uc := NewListVehicles(repo)
	got, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 vehicle, got %d", len(got))
	}
}

func TestListVehicles_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	repo := mocks.NewMockVehicleRepository(ctrl)
	repo.EXPECT().FindAll(gomock.Any()).Return(nil, repoErr)

	uc := NewListVehicles(repo)
	_, err := uc.Execute(context.Background())
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
