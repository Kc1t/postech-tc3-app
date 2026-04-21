package serviceorderuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestDeleteServiceOrder_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	so := entities.NewServiceOrder("cust-1", "veh-1") // status = received
	so.SetID("so-1")

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "so-1").Return(so, nil)
	repo.EXPECT().Delete(gomock.Any(), "so-1").Return(nil)

	uc := NewDeleteServiceOrder(repo)
	if err := uc.Execute(context.Background(), "so-1"); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestDeleteServiceOrder_Execute_NotCancellable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	so := entities.NewServiceOrder("cust-1", "veh-1")
	so.SetID("so-1")
	_ = so.UpdateStatus(entities.StatusInDiagnosis) // mover para status nao cancelavel

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "so-1").Return(so, nil)

	uc := NewDeleteServiceOrder(repo)
	err := uc.Execute(context.Background(), "so-1")
	if !errors.Is(err, domainerrors.ErrOrderNotCancellable) {
		t.Fatalf("expected ErrOrderNotCancellable, got %v", err)
	}
}

func TestDeleteServiceOrder_Execute_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "so-x").Return(nil, domainerrors.ErrNotFound)

	uc := NewDeleteServiceOrder(repo)
	err := uc.Execute(context.Background(), "so-x")
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteServiceOrder_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "so-1").Return(nil, repoErr)

	uc := NewDeleteServiceOrder(repo)
	if err := uc.Execute(context.Background(), "so-1"); !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
