package serviceorderuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestCreateServiceOrder_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	customerRepo := mocks.NewMockCustomerRepository(ctrl)
	vehicleRepo := mocks.NewMockVehicleRepository(ctrl)
	serviceRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	soRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	uc := NewCreateServiceOrder(soRepo, customerRepo, vehicleRepo, serviceRepo, partRepo)
	so := entities.NewServiceOrder("cust-1", "veh-1")
	if err := uc.Execute(context.Background(), so); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if so.Status() != entities.StatusReceived {
		t.Errorf("expected status %q, got %q", entities.StatusReceived, so.Status())
	}
}

func TestCreateServiceOrder_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	customerRepo := mocks.NewMockCustomerRepository(ctrl)
	vehicleRepo := mocks.NewMockVehicleRepository(ctrl)
	serviceRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	soRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(repoErr)

	uc := NewCreateServiceOrder(soRepo, customerRepo, vehicleRepo, serviceRepo, partRepo)
	so := entities.NewServiceOrder("cust-1", "veh-1")
	if err := uc.Execute(context.Background(), so); !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
