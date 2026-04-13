package vehicleuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestUpdateVehicle_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	v := entities.NewVehicle("cust-1", "ABC1234", "Toyota", "Corolla", 2020)
	repo := mocks.NewMockVehicleRepository(ctrl)
	repo.EXPECT().Update(gomock.Any(), v).Return(nil)

	uc := NewUpdateVehicle(repo)
	if err := uc.Execute(context.Background(), v); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestUpdateVehicle_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	v := entities.NewVehicle("cust-1", "ABC1234", "Toyota", "Corolla", 2020)
	repo := mocks.NewMockVehicleRepository(ctrl)
	repo.EXPECT().Update(gomock.Any(), v).Return(repoErr)

	uc := NewUpdateVehicle(repo)
	if err := uc.Execute(context.Background(), v); !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
