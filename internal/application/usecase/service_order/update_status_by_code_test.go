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

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	so := entities.ReconstituteServiceOrder("order-1", 100, "cust-1", "veh-1",
		entities.StatusAwaitingApproval, nil, nil, 0, "", ft(), ft(), nil, nil)

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil)
	soRepo.EXPECT().FindByCode(gomock.Any(), 100).Return(so, nil)
	// Aprovacao sem pecas: nao chama UpdateStock; Update (nao UpdateStatus)
	// porque startedAt acabou de ser gravado pela entidade.
	soRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, partRepo)
	if err := uc.Execute(context.Background(), 100, "52998224725", entities.StatusInExecution); err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
}

func TestUpdateStatusByCode_Rejeicao_Sucesso(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	so := entities.ReconstituteServiceOrder("order-1", 100, "cust-1", "veh-1",
		entities.StatusAwaitingApproval, nil, nil, 0, "", ft(), ft(), nil, nil)

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil)
	soRepo.EXPECT().FindByCode(gomock.Any(), 100).Return(so, nil)
	soRepo.EXPECT().UpdateStatus(gomock.Any(), "order-1", entities.StatusReceived).Return(nil)

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, partRepo)
	if err := uc.Execute(context.Background(), 100, "52998224725", entities.StatusReceived); err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
}

func TestUpdateStatusByCode_ClienteNaoEncontrado(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(nil, domainerrors.ErrNotFound)

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, partRepo)
	err := uc.Execute(context.Background(), 100, "52998224725", entities.StatusInExecution)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("erro = %v, esperava ErrNotFound", err)
	}
}

func TestUpdateStatusByCode_OSNaoEncontrada(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil)
	soRepo.EXPECT().FindByCode(gomock.Any(), 999).Return(nil, domainerrors.ErrNotFound)

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, partRepo)
	err := uc.Execute(context.Background(), 999, "52998224725", entities.StatusInExecution)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("erro = %v, esperava ErrNotFound", err)
	}
}

func TestUpdateStatusByCode_OSDeOutroCliente(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	so := entities.ReconstituteServiceOrder("order-1", 100, "outro-cliente", "veh-1",
		entities.StatusAwaitingApproval, nil, nil, 0, "", ft(), ft(), nil, nil)

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil)
	soRepo.EXPECT().FindByCode(gomock.Any(), 100).Return(so, nil)

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, partRepo)
	err := uc.Execute(context.Background(), 100, "52998224725", entities.StatusInExecution)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("erro = %v, esperava ErrNotFound (ownership)", err)
	}
}

func TestUpdateStatusByCode_TransicaoInvalida(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	so := entities.ReconstituteServiceOrder("order-1", 100, "cust-1", "veh-1",
		entities.StatusReceived, nil, nil, 0, "", ft(), ft(), nil, nil)

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil)
	soRepo.EXPECT().FindByCode(gomock.Any(), 100).Return(so, nil)

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, partRepo)
	err := uc.Execute(context.Background(), 100, "52998224725", entities.StatusInExecution)
	if !errors.Is(err, domainerrors.ErrInvalidStatus) {
		t.Fatalf("erro = %v, esperava ErrInvalidStatus", err)
	}
}

func TestUpdateStatusByCode_StatusNaoPermitidoParaCliente(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	// OS em in_execution; cliente tenta avancar para finished (acao do mecanico).
	so := entities.ReconstituteServiceOrder("order-1", 100, "cust-1", "veh-1",
		entities.StatusInExecution, nil, nil, 0, "", ft(), ft(), nil, nil)

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil)
	soRepo.EXPECT().FindByCode(gomock.Any(), 100).Return(so, nil)
	// Nao deve chamar UpdateStatus no repositorio.

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, partRepo)
	err := uc.Execute(context.Background(), 100, "52998224725", entities.StatusFinished)
	if !errors.Is(err, domainerrors.ErrStatusNotAllowedForCustomer) {
		t.Fatalf("erro = %v, esperava ErrStatusNotAllowedForCustomer", err)
	}
}

func TestUpdateStatusByCode_ErroInfraCliente(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	infraErr := errors.New("timeout")
	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(nil, infraErr)

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, partRepo)
	err := uc.Execute(context.Background(), 100, "52998224725", entities.StatusInExecution)
	if !errors.Is(err, infraErr) {
		t.Fatalf("erro = %v, esperava erro de infra propagado", err)
	}
}

// Aprovacao pelo cliente (rota publica) tambem deve baixar o estoque das
// pecas do orcamento antes de transitar a OS.
func TestUpdateStatusByCode_Aprovacao_DecrementaEstoque(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	parts := []entities.PartItem{
		{PartID: "part-1", Description: "Filtro", Quantity: 2, UnitPrice: 25.00},
		{PartID: "part-2", Description: "Oleo", Quantity: 4, UnitPrice: 40.00},
	}
	so := entities.ReconstituteServiceOrder("order-1", 100, "cust-1", "veh-1",
		entities.StatusAwaitingApproval, nil, parts, 0, "", ft(), ft(), nil, nil)

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	gomock.InOrder(
		custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil),
		soRepo.EXPECT().FindByCode(gomock.Any(), 100).Return(so, nil),
		partRepo.EXPECT().UpdateStock(gomock.Any(), "part-1", -2).Return(nil),
		partRepo.EXPECT().UpdateStock(gomock.Any(), "part-2", -4).Return(nil),
		// Update (nao UpdateStatus) porque startedAt foi gravado pela entidade.
		soRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, saved *entities.ServiceOrder) error {
			if saved.StartedAt() == nil {
				t.Error("StartedAt() = nil, esperava timestamp gravado na aprovacao")
			}
			return nil
		}),
	)

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, partRepo)
	if err := uc.Execute(context.Background(), 100, "52998224725", entities.StatusInExecution); err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
}

// Recusa do orcamento pelo cliente NAO deve mexer em estoque.
func TestUpdateStatusByCode_Rejeicao_NaoTocaEstoque(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	parts := []entities.PartItem{
		{PartID: "part-1", Description: "Filtro", Quantity: 2, UnitPrice: 25.00},
	}
	so := entities.ReconstituteServiceOrder("order-1", 100, "cust-1", "veh-1",
		entities.StatusAwaitingApproval, nil, parts, 0, "", ft(), ft(), nil, nil)

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil)
	soRepo.EXPECT().FindByCode(gomock.Any(), 100).Return(so, nil)
	soRepo.EXPECT().UpdateStatus(gomock.Any(), "order-1", entities.StatusReceived).Return(nil)
	// partRepo.UpdateStock NAO deve ser chamado.

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, partRepo)
	if err := uc.Execute(context.Background(), 100, "52998224725", entities.StatusReceived); err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
}

// Estoque insuficiente no momento da aprovacao do cliente: peca ja
// decrementada e revertida, status nao muda.
func TestUpdateStatusByCode_Aprovacao_EstoqueInsuficiente_FazRollback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	parts := []entities.PartItem{
		{PartID: "part-1", Description: "Filtro", Quantity: 2, UnitPrice: 25.00},
		{PartID: "part-2", Description: "Oleo", Quantity: 10, UnitPrice: 40.00},
	}
	so := entities.ReconstituteServiceOrder("order-1", 100, "cust-1", "veh-1",
		entities.StatusAwaitingApproval, nil, parts, 0, "", ft(), ft(), nil, nil)

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	gomock.InOrder(
		custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil),
		soRepo.EXPECT().FindByCode(gomock.Any(), 100).Return(so, nil),
		partRepo.EXPECT().UpdateStock(gomock.Any(), "part-1", -2).Return(nil),
		partRepo.EXPECT().UpdateStock(gomock.Any(), "part-2", -10).Return(domainerrors.ErrInsufficientStock),
		partRepo.EXPECT().UpdateStock(gomock.Any(), "part-1", 2).Return(nil),
	)
	// soRepo.UpdateStatus NAO deve ser chamado.

	uc := NewUpdateServiceOrderStatusByCode(soRepo, custRepo, partRepo)
	err := uc.Execute(context.Background(), 100, "52998224725", entities.StatusInExecution)
	if !errors.Is(err, domainerrors.ErrInsufficientStock) {
		t.Fatalf("erro = %v, esperava ErrInsufficientStock", err)
	}
}
