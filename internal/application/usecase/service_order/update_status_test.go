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

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusReceived, nil, nil, 0, "", ft(), ft(), nil, nil)

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	repo.EXPECT().UpdateStatus(gomock.Any(), "order-1", entities.StatusInDiagnosis).Return(nil)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := entities.StatusUpdate{
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

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusReceived, nil, nil, 0, "", ft(), ft(), nil, nil)

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	// UpdateStatus NAO deve ser chamado

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := entities.StatusUpdate{
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
	input := entities.StatusUpdate{
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

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusReceived, nil, nil, 0, "", ft(), ft(), nil, nil)

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	// UpdateStatus NAO deve ser chamado

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := entities.StatusUpdate{
		ID:     "order-1",
		Status: entities.OrderStatus("invalido"),
	}
	err := uc.Execute(context.Background(), input)
	if err == nil {
		t.Fatal("esperava erro para status desconhecido")
	}
}

// =============================================================================
// Testes do branch awaiting_approval (montagem de orcamento)
// =============================================================================

func TestUpdateStatus_AwaitingApproval_ComServicosEPecas(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusInDiagnosis, nil, nil, 0, "", ft(), ft(), nil, nil)

	svc1 := entities.ReconstituteService("svc-1", 1, "Troca de oleo", "Desc", 150.00, 30, ft(), ft())
	svc2 := entities.ReconstituteService("svc-2", 2, "Alinhamento", "Desc", 80.00, 45, ft(), ft())
	part1 := entities.ReconstitutePart("part-1", "FAB-001", "Filtro", "Desc", "un", 25.00, 50, ft(), ft())

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	svcRepo.EXPECT().FindByCodes(gomock.Any(), []int{1, 2}).Return([]*entities.Service{svc1, svc2}, nil)
	partRepo.EXPECT().FindByManufacturerCodes(gomock.Any(), []string{"FAB-001"}).Return([]*entities.Part{part1}, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := entities.StatusUpdate{
		ID:           "order-1",
		Status:       entities.StatusAwaitingApproval,
		ServiceCodes: []int{1, 2},
		Parts: []entities.OrderPartItem{
			{ManufacturerCode: "FAB-001", Quantity: 2},
		},
	}
	if err := uc.Execute(context.Background(), input); err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}

	// 150 + 80 + (2 * 25) = 280
	if so.TotalAmount() != 280.00 {
		t.Errorf("TotalAmount = %.2f, esperava 280.00", so.TotalAmount())
	}
	if len(so.Services()) != 2 {
		t.Errorf("len(Services) = %d, esperava 2", len(so.Services()))
	}
	if len(so.Parts()) != 1 {
		t.Errorf("len(Parts) = %d, esperava 1", len(so.Parts()))
	}
}

func TestUpdateStatus_AwaitingApproval_SemServicosNemPecas(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusInDiagnosis, nil, nil, 0, "", ft(), ft(), nil, nil)

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := entities.StatusUpdate{
		ID:     "order-1",
		Status: entities.StatusAwaitingApproval,
	}
	if err := uc.Execute(context.Background(), input); err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
}

func TestUpdateStatus_AwaitingApproval_ServicoInexistente(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusInDiagnosis, nil, nil, 0, "", ft(), ft(), nil, nil)

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	// Retorna lista vazia — codigo 999 nao existe
	svcRepo.EXPECT().FindByCodes(gomock.Any(), []int{999}).Return([]*entities.Service{}, nil)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := entities.StatusUpdate{
		ID:           "order-1",
		Status:       entities.StatusAwaitingApproval,
		ServiceCodes: []int{999},
	}
	err := uc.Execute(context.Background(), input)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("erro = %v, esperava ErrNotFound (servico inexistente)", err)
	}
}

func TestUpdateStatus_AwaitingApproval_PecaInexistente(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusInDiagnosis, nil, nil, 0, "", ft(), ft(), nil, nil)

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	// Retorna lista vazia — peca nao existe
	partRepo.EXPECT().FindByManufacturerCodes(gomock.Any(), []string{"INEXISTENTE"}).Return([]*entities.Part{}, nil)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := entities.StatusUpdate{
		ID:     "order-1",
		Status: entities.StatusAwaitingApproval,
		Parts: []entities.OrderPartItem{
			{ManufacturerCode: "INEXISTENTE", Quantity: 1},
		},
	}
	err := uc.Execute(context.Background(), input)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("erro = %v, esperava ErrNotFound (peca inexistente)", err)
	}
}

func TestUpdateStatus_AwaitingApproval_EstoqueInsuficiente(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusInDiagnosis, nil, nil, 0, "", ft(), ft(), nil, nil)
	part := entities.ReconstitutePart("part-1", "FAB-001", "Filtro", "Desc", "un", 25.00, 2, ft(), ft())

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	partRepo.EXPECT().FindByManufacturerCodes(gomock.Any(), []string{"FAB-001"}).Return([]*entities.Part{part}, nil)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := entities.StatusUpdate{
		ID:     "order-1",
		Status: entities.StatusAwaitingApproval,
		Parts: []entities.OrderPartItem{
			{ManufacturerCode: "FAB-001", Quantity: 10}, // pede 10, tem 2
		},
	}
	err := uc.Execute(context.Background(), input)
	if !errors.Is(err, domainerrors.ErrInsufficientStock) {
		t.Fatalf("erro = %v, esperava ErrInsufficientStock", err)
	}
}

func TestUpdateStatus_AwaitingApproval_ErroInfraServiceRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusInDiagnosis, nil, nil, 0, "", ft(), ft(), nil, nil)
	infraErr := errors.New("db timeout")

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	svcRepo.EXPECT().FindByCodes(gomock.Any(), []int{1}).Return(nil, infraErr)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := entities.StatusUpdate{
		ID:           "order-1",
		Status:       entities.StatusAwaitingApproval,
		ServiceCodes: []int{1},
	}
	err := uc.Execute(context.Background(), input)
	if !errors.Is(err, infraErr) {
		t.Fatalf("erro = %v, esperava erro de infra propagado", err)
	}
}

func TestUpdateStatus_AwaitingApproval_ErroInfraPartRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusInDiagnosis, nil, nil, 0, "", ft(), ft(), nil, nil)
	infraErr := errors.New("db timeout")

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	partRepo.EXPECT().FindByManufacturerCodes(gomock.Any(), []string{"FAB-001"}).Return(nil, infraErr)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := entities.StatusUpdate{
		ID:     "order-1",
		Status: entities.StatusAwaitingApproval,
		Parts: []entities.OrderPartItem{
			{ManufacturerCode: "FAB-001", Quantity: 1},
		},
	}
	err := uc.Execute(context.Background(), input)
	if !errors.Is(err, infraErr) {
		t.Fatalf("erro = %v, esperava erro de infra propagado", err)
	}
}

func TestUpdateStatus_AwaitingApproval_ErroRepoUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusInDiagnosis, nil, nil, 0, "", ft(), ft(), nil, nil)
	infraErr := errors.New("db write error")

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(infraErr)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := entities.StatusUpdate{
		ID:     "order-1",
		Status: entities.StatusAwaitingApproval,
	}
	err := uc.Execute(context.Background(), input)
	if !errors.Is(err, infraErr) {
		t.Fatalf("erro = %v, esperava erro de infra propagado", err)
	}
}

func TestUpdateStatus_ErroRepoUpdateStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusReceived, nil, nil, 0, "", ft(), ft(), nil, nil)
	infraErr := errors.New("db write error")

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	repo.EXPECT().UpdateStatus(gomock.Any(), "order-1", entities.StatusInDiagnosis).Return(infraErr)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := entities.StatusUpdate{
		ID:     "order-1",
		Status: entities.StatusInDiagnosis,
	}
	err := uc.Execute(context.Background(), input)
	if !errors.Is(err, infraErr) {
		t.Fatalf("erro = %v, esperava erro de infra propagado", err)
	}
}

// =============================================================================
// Testes da aprovacao (awaiting_approval -> in_execution) com baixa de estoque
// =============================================================================

func TestUpdateStatus_InExecution_DecrementaEstoqueETransita(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	parts := []entities.PartItem{
		{PartID: "part-1", Description: "Filtro", Quantity: 2, UnitPrice: 25.00},
		{PartID: "part-2", Description: "Oleo", Quantity: 4, UnitPrice: 40.00},
	}
	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1",
		entities.StatusAwaitingApproval, nil, parts, 0, "", ft(), ft(), nil, nil)

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	gomock.InOrder(
		repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil),
		partRepo.EXPECT().UpdateStock(gomock.Any(), "part-1", -2).Return(nil),
		partRepo.EXPECT().UpdateStock(gomock.Any(), "part-2", -4).Return(nil),
		// Update (nao UpdateStatus) porque a entidade acabou de gravar startedAt.
		repo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, saved *entities.ServiceOrder) error {
			if saved.Status() != entities.StatusInExecution {
				t.Errorf("Status persistido = %q, esperava in_execution", saved.Status())
			}
			if saved.StartedAt() == nil {
				t.Error("StartedAt() = nil, esperava timestamp gravado na transicao")
			}
			return nil
		}),
	)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := entities.StatusUpdate{ID: "order-1", Status: entities.StatusInExecution}
	if err := uc.Execute(context.Background(), input); err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
}

func TestUpdateStatus_InExecution_SemPecas_NaoTocaEstoque(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// OS sem pecas (apenas servicos, por exemplo).
	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1",
		entities.StatusAwaitingApproval, nil, nil, 0, "", ft(), ft(), nil, nil)

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	// partRepo.UpdateStock NAO deve ser chamado
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := entities.StatusUpdate{ID: "order-1", Status: entities.StatusInExecution}
	if err := uc.Execute(context.Background(), input); err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
}

// Se a segunda peca falhar, a primeira (ja decrementada) deve ser revertida
// e a transicao de status NAO pode ocorrer.
func TestUpdateStatus_InExecution_EstoqueInsuficiente_FazRollback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	parts := []entities.PartItem{
		{PartID: "part-1", Description: "Filtro", Quantity: 2, UnitPrice: 25.00},
		{PartID: "part-2", Description: "Oleo", Quantity: 5, UnitPrice: 40.00},
		{PartID: "part-3", Description: "Pastilha", Quantity: 1, UnitPrice: 180.00},
	}
	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1",
		entities.StatusAwaitingApproval, nil, parts, 0, "", ft(), ft(), nil, nil)

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	gomock.InOrder(
		repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil),
		partRepo.EXPECT().UpdateStock(gomock.Any(), "part-1", -2).Return(nil),
		partRepo.EXPECT().UpdateStock(gomock.Any(), "part-2", -5).Return(domainerrors.ErrInsufficientStock),
		// rollback da unica peca ja decrementada (part-1)
		partRepo.EXPECT().UpdateStock(gomock.Any(), "part-1", 2).Return(nil),
	)
	// repo.UpdateStatus NAO deve ser chamado — transicao nao ocorre.
	// UpdateStock de part-3 NAO deve ser chamado — fail-fast.

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := entities.StatusUpdate{ID: "order-1", Status: entities.StatusInExecution}
	err := uc.Execute(context.Background(), input)
	if !errors.Is(err, domainerrors.ErrInsufficientStock) {
		t.Fatalf("erro = %v, esperava ErrInsufficientStock", err)
	}
}

// Erro no rollback nao mascara o erro original — o use case retorna o erro
// que causou a falha, nao o erro da reversao.
func TestUpdateStatus_InExecution_FalhaNoRollback_PropagaErroOriginal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	parts := []entities.PartItem{
		{PartID: "part-1", Description: "Filtro", Quantity: 2, UnitPrice: 25.00},
		{PartID: "part-2", Description: "Oleo", Quantity: 5, UnitPrice: 40.00},
	}
	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1",
		entities.StatusAwaitingApproval, nil, parts, 0, "", ft(), ft(), nil, nil)

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	rollbackErr := errors.New("db unavailable during rollback")
	gomock.InOrder(
		repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil),
		partRepo.EXPECT().UpdateStock(gomock.Any(), "part-1", -2).Return(nil),
		partRepo.EXPECT().UpdateStock(gomock.Any(), "part-2", -5).Return(domainerrors.ErrInsufficientStock),
		partRepo.EXPECT().UpdateStock(gomock.Any(), "part-1", 2).Return(rollbackErr),
	)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := entities.StatusUpdate{ID: "order-1", Status: entities.StatusInExecution}
	err := uc.Execute(context.Background(), input)
	if !errors.Is(err, domainerrors.ErrInsufficientStock) {
		t.Fatalf("erro = %v, esperava ErrInsufficientStock (erro original, nao %v)", err, rollbackErr)
	}
}
