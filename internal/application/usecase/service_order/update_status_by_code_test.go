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

func TestUpdateStatusByCode_Aprovacao_Sucesso(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requester := entities.ReconstituteRequester("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	so := entities.ReconstituteServiceOrder("order-1", 100, "cust-1", "veh-1",
		entities.StatusAwaitingApproval, nil, nil, 0, "", ft(), ft(), nil, nil)

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockRequesterRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(requester, nil)
	soRepo.EXPECT().FindByCode(gomock.Any(), 100).Return(so, nil)
	// Aprovacao: use case delega persistencia + decremento atomicos ao repo.
	soRepo.EXPECT().ApplyApprovalTransition(gomock.Any(), gomock.Any()).Return(nil)

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, nil)
	if err := uc.Execute(context.Background(), 100, "52998224725", entities.StatusInExecution); err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
}

func TestUpdateStatusByCode_Rejeicao_Sucesso(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requester := entities.ReconstituteRequester("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	so := entities.ReconstituteServiceOrder("order-1", 100, "cust-1", "veh-1",
		entities.StatusAwaitingApproval, nil, nil, 0, "", ft(), ft(), nil, nil)

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockRequesterRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(requester, nil)
	soRepo.EXPECT().FindByCode(gomock.Any(), 100).Return(so, nil)
	soRepo.EXPECT().UpdateStatus(gomock.Any(), "order-1", entities.StatusReceived).Return(nil)

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, nil)
	if err := uc.Execute(context.Background(), 100, "52998224725", entities.StatusReceived); err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
}

func TestUpdateStatusByCode_ClienteNaoEncontrado(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockRequesterRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(nil, domainerrors.ErrNotFound)

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, nil)
	err := uc.Execute(context.Background(), 100, "52998224725", entities.StatusInExecution)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("erro = %v, esperava ErrNotFound", err)
	}
}

func TestUpdateStatusByCode_OSNaoEncontrada(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requester := entities.ReconstituteRequester("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockRequesterRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(requester, nil)
	soRepo.EXPECT().FindByCode(gomock.Any(), 999).Return(nil, domainerrors.ErrNotFound)

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, nil)
	err := uc.Execute(context.Background(), 999, "52998224725", entities.StatusInExecution)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("erro = %v, esperava ErrNotFound", err)
	}
}

func TestUpdateStatusByCode_OSDeOutroCliente(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requester := entities.ReconstituteRequester("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	so := entities.ReconstituteServiceOrder("order-1", 100, "outro-cliente", "veh-1",
		entities.StatusAwaitingApproval, nil, nil, 0, "", ft(), ft(), nil, nil)

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockRequesterRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(requester, nil)
	soRepo.EXPECT().FindByCode(gomock.Any(), 100).Return(so, nil)

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, nil)
	err := uc.Execute(context.Background(), 100, "52998224725", entities.StatusInExecution)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("erro = %v, esperava ErrNotFound (ownership)", err)
	}
}

func TestUpdateStatusByCode_TransicaoInvalida(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requester := entities.ReconstituteRequester("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	so := entities.ReconstituteServiceOrder("order-1", 100, "cust-1", "veh-1",
		entities.StatusReceived, nil, nil, 0, "", ft(), ft(), nil, nil)

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockRequesterRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(requester, nil)
	soRepo.EXPECT().FindByCode(gomock.Any(), 100).Return(so, nil)

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, nil)
	err := uc.Execute(context.Background(), 100, "52998224725", entities.StatusInExecution)
	if !errors.Is(err, domainerrors.ErrInvalidStatus) {
		t.Fatalf("erro = %v, esperava ErrInvalidStatus", err)
	}
}

func TestUpdateStatusByCode_StatusNaoPermitidoParaCliente(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requester := entities.ReconstituteRequester("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	// OS em in_execution; cliente tenta avancar para finished (acao do mecanico).
	so := entities.ReconstituteServiceOrder("order-1", 100, "cust-1", "veh-1",
		entities.StatusInExecution, nil, nil, 0, "", ft(), ft(), nil, nil)

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockRequesterRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(requester, nil)
	soRepo.EXPECT().FindByCode(gomock.Any(), 100).Return(so, nil)
	// Nao deve chamar UpdateStatus no repositorio.

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, nil)
	err := uc.Execute(context.Background(), 100, "52998224725", entities.StatusFinished)
	if !errors.Is(err, domainerrors.ErrStatusNotAllowedForRequester) {
		t.Fatalf("erro = %v, esperava ErrStatusNotAllowedForRequester", err)
	}
}

func TestUpdateStatusByCode_ErroInfraCliente(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	infraErr := errors.New("timeout")
	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockRequesterRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(nil, infraErr)

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, nil)
	err := uc.Execute(context.Background(), 100, "52998224725", entities.StatusInExecution)
	if !errors.Is(err, infraErr) {
		t.Fatalf("erro = %v, esperava erro de infra propagado", err)
	}
}

// Aprovacao pelo cliente delega para ApplyApprovalTransition — o use case
// so precisa garantir que a OS passa em estado in_execution com startedAt
// gravado pela entidade antes de chegar ao repo.
func TestUpdateStatusByCode_Aprovacao_DelegaParaApplyApprovalTransition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requester := entities.ReconstituteRequester("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	parts := []entities.PartItem{
		{PartID: "part-1", Description: "Filtro", Quantity: 2, UnitPrice: 25.00},
		{PartID: "part-2", Description: "Oleo", Quantity: 4, UnitPrice: 40.00},
	}
	so := entities.ReconstituteServiceOrder("order-1", 100, "cust-1", "veh-1",
		entities.StatusAwaitingApproval, nil, parts, 0, "", ft(), ft(), nil, nil)

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockRequesterRepository(ctrl)

	gomock.InOrder(
		custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(requester, nil),
		soRepo.EXPECT().FindByCode(gomock.Any(), 100).Return(so, nil),
		soRepo.EXPECT().ApplyApprovalTransition(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, saved *entities.ServiceOrder) error {
			if saved.Status() != entities.StatusInExecution {
				t.Errorf("Status persistido = %q, esperava in_execution", saved.Status())
			}
			if saved.StartedAt() == nil {
				t.Error("StartedAt() = nil, esperava timestamp gravado na aprovacao")
			}
			return nil
		}),
	)

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, nil)
	if err := uc.Execute(context.Background(), 100, "52998224725", entities.StatusInExecution); err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
}

// Recusa do orcamento pelo cliente NAO deve invocar o caminho transacional.
func TestUpdateStatusByCode_Rejeicao_NaoInvocaApplyApprovalTransition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requester := entities.ReconstituteRequester("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	parts := []entities.PartItem{
		{PartID: "part-1", Description: "Filtro", Quantity: 2, UnitPrice: 25.00},
	}
	so := entities.ReconstituteServiceOrder("order-1", 100, "cust-1", "veh-1",
		entities.StatusAwaitingApproval, nil, parts, 0, "", ft(), ft(), nil, nil)

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockRequesterRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(requester, nil)
	soRepo.EXPECT().FindByCode(gomock.Any(), 100).Return(so, nil)
	soRepo.EXPECT().UpdateStatus(gomock.Any(), "order-1", entities.StatusReceived).Return(nil)
	// ApplyApprovalTransition NAO deve ser chamado.

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, nil)
	if err := uc.Execute(context.Background(), 100, "52998224725", entities.StatusReceived); err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
}

// Estoque insuficiente vem como erro do repo (rollback ja acontece na tx DB).
func TestUpdateStatusByCode_Aprovacao_PropagaErroDoRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requester := entities.ReconstituteRequester("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	parts := []entities.PartItem{
		{PartID: "part-1", Description: "Filtro", Quantity: 2, UnitPrice: 25.00},
	}
	so := entities.ReconstituteServiceOrder("order-1", 100, "cust-1", "veh-1",
		entities.StatusAwaitingApproval, nil, parts, 0, "", ft(), ft(), nil, nil)

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockRequesterRepository(ctrl)

	gomock.InOrder(
		custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(requester, nil),
		soRepo.EXPECT().FindByCode(gomock.Any(), 100).Return(so, nil),
		soRepo.EXPECT().ApplyApprovalTransition(gomock.Any(), gomock.Any()).Return(domainerrors.ErrInsufficientStock),
	)

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, nil)
	err := uc.Execute(context.Background(), 100, "52998224725", entities.StatusInExecution)
	if !errors.Is(err, domainerrors.ErrInsufficientStock) {
		t.Fatalf("erro = %v, esperava ErrInsufficientStock", err)
	}
}
