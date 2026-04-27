package serviceorderuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestListServiceOrdersByRequester_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := []*entities.ServiceOrder{
		entities.NewServiceOrder("cust-1", "veh-1"),
	}
	repo := mocks.NewMockServiceOrderRepository(ctrl)
	repo.EXPECT().FindByRequesterID(gomock.Any(), "cust-1").Return(expected, nil)

	uc := NewListServiceOrdersByRequester(repo)
	got, err := uc.Execute(context.Background(), "cust-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 service order, got %d", len(got))
	}
}

func TestListServiceOrdersByRequester_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	repo := mocks.NewMockServiceOrderRepository(ctrl)
	repo.EXPECT().FindByRequesterID(gomock.Any(), "cust-1").Return(nil, repoErr)

	uc := NewListServiceOrdersByRequester(repo)
	_, err := uc.Execute(context.Background(), "cust-1")
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}
