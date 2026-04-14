package serviceorderuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestUpdateServiceOrder_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	so := entities.NewServiceOrder("cust-1", "veh-1")
	repo := mocks.NewMockServiceOrderRepository(ctrl)
	repo.EXPECT().Update(gomock.Any(), so).Return(nil)

	uc := NewUpdateServiceOrder(repo)
	if err := uc.Execute(context.Background(), so); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestUpdateServiceOrder_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	so := entities.NewServiceOrder("cust-1", "veh-1")
	repo := mocks.NewMockServiceOrderRepository(ctrl)
	repo.EXPECT().Update(gomock.Any(), so).Return(repoErr)

	uc := NewUpdateServiceOrder(repo)
	if err := uc.Execute(context.Background(), so); !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
