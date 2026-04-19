package vehicleuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestDeleteVehicle_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockVehicleRepository(ctrl)
	repo.EXPECT().Delete(gomock.Any(), "veh-1").Return(nil)

	uc := NewDeleteVehicle(repo)
	if err := uc.Execute(context.Background(), "veh-1"); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestDeleteVehicle_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	repo := mocks.NewMockVehicleRepository(ctrl)
	repo.EXPECT().Delete(gomock.Any(), "veh-1").Return(repoErr)

	uc := NewDeleteVehicle(repo)
	if err := uc.Execute(context.Background(), "veh-1"); !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
