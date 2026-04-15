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

func TestUpdateStatus_TransicaoValida(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	so := entities.ReconstituteServiceOrder("order-1", "c1", "v1", entities.StatusReceived, nil, nil, 0, "", ft(), ft())

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	repo.EXPECT().UpdateStatus(gomock.Any(), "order-1", entities.StatusInDiagnosis).Return(nil)

	uc := NewUpdateServiceOrderStatus(repo)
	if err := uc.Execute(context.Background(), "order-1", entities.StatusInDiagnosis); err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
}

func TestUpdateStatus_TransicaoInvalida(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	so := entities.ReconstituteServiceOrder("order-1", "c1", "v1", entities.StatusReceived, nil, nil, 0, "", ft(), ft())

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	// UpdateStatus NAO deve ser chamado

	uc := NewUpdateServiceOrderStatus(repo)
	err := uc.Execute(context.Background(), "order-1", entities.StatusFinished)
	if !errors.Is(err, domainerrors.ErrInvalidStatus) {
		t.Fatalf("erro = %v, esperava ErrInvalidStatus", err)
	}
}

func TestUpdateStatus_OSNaoEncontrada(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "inexistente").Return(nil, domainerrors.ErrNotFound)

	uc := NewUpdateServiceOrderStatus(repo)
	err := uc.Execute(context.Background(), "inexistente", entities.StatusInDiagnosis)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("erro = %v, esperava ErrNotFound", err)
	}
}

func TestUpdateStatus_StatusDesconhecido(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	so := entities.ReconstituteServiceOrder("order-1", "c1", "v1", entities.StatusReceived, nil, nil, 0, "", ft(), ft())

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	// UpdateStatus NAO deve ser chamado

	uc := NewUpdateServiceOrderStatus(repo)
	err := uc.Execute(context.Background(), "order-1", entities.OrderStatus("invalido"))
	if err == nil {
		t.Fatal("esperava erro para status desconhecido")
	}
}
