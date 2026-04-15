package serviceorderuc

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestUpdateStatus_TransicaoValida(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusReceived, nil, nil, 0, "", ft(), ft())

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	repo.EXPECT().UpdateStatus(gomock.Any(), "order-1", entities.StatusInDiagnosis).Return(nil)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := ports.UpdateStatusInput{
		ID:     "order-1",
		Status: entities.StatusInDiagnosis,
	}
	if err := uc.Execute(context.Background(), input); err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
}

func TestUpdateStatus_TransicaoInvalida(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusReceived, nil, nil, 0, "", ft(), ft())

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	// UpdateStatus NAO deve ser chamado

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := ports.UpdateStatusInput{
		ID:     "order-1",
		Status: entities.StatusFinished,
	}
	err := uc.Execute(context.Background(), input)
	if !errors.Is(err, domainerrors.ErrInvalidStatus) {
		t.Fatalf("erro = %v, esperava ErrInvalidStatus", err)
	}
}

func TestUpdateStatus_OSNaoEncontrada(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "inexistente").Return(nil, domainerrors.ErrNotFound)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := ports.UpdateStatusInput{
		ID:     "inexistente",
		Status: entities.StatusInDiagnosis,
	}
	err := uc.Execute(context.Background(), input)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("erro = %v, esperava ErrNotFound", err)
	}
}

func TestUpdateStatus_StatusDesconhecido(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusReceived, nil, nil, 0, "", ft(), ft())

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	// UpdateStatus NAO deve ser chamado

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := ports.UpdateStatusInput{
		ID:     "order-1",
		Status: entities.OrderStatus("invalido"),
	}
	err := uc.Execute(context.Background(), input)
	if err == nil {
		t.Fatal("esperava erro para status desconhecido")
	}
}
