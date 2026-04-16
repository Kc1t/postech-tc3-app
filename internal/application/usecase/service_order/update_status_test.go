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

// =============================================================================
// Testes do branch awaiting_approval (montagem de orcamento)
// =============================================================================

func TestUpdateStatus_AwaitingApproval_ComServicosEPecas(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusInDiagnosis, nil, nil, 0, "", ft(), ft())

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
	input := ports.UpdateStatusInput{
		ID:           "order-1",
		Status:       entities.StatusAwaitingApproval,
		ServiceCodes: []int{1, 2},
		Parts: []ports.UpdateStatusPartInput{
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

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusInDiagnosis, nil, nil, 0, "", ft(), ft())

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := ports.UpdateStatusInput{
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

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusInDiagnosis, nil, nil, 0, "", ft(), ft())

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	// Retorna lista vazia — codigo 999 nao existe
	svcRepo.EXPECT().FindByCodes(gomock.Any(), []int{999}).Return([]*entities.Service{}, nil)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := ports.UpdateStatusInput{
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

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusInDiagnosis, nil, nil, 0, "", ft(), ft())

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	// Retorna lista vazia — peca nao existe
	partRepo.EXPECT().FindByManufacturerCodes(gomock.Any(), []string{"INEXISTENTE"}).Return([]*entities.Part{}, nil)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := ports.UpdateStatusInput{
		ID:     "order-1",
		Status: entities.StatusAwaitingApproval,
		Parts: []ports.UpdateStatusPartInput{
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

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusInDiagnosis, nil, nil, 0, "", ft(), ft())
	part := entities.ReconstitutePart("part-1", "FAB-001", "Filtro", "Desc", "un", 25.00, 2, ft(), ft())

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	partRepo.EXPECT().FindByManufacturerCodes(gomock.Any(), []string{"FAB-001"}).Return([]*entities.Part{part}, nil)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := ports.UpdateStatusInput{
		ID:     "order-1",
		Status: entities.StatusAwaitingApproval,
		Parts: []ports.UpdateStatusPartInput{
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

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusInDiagnosis, nil, nil, 0, "", ft(), ft())
	infraErr := errors.New("db timeout")

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	svcRepo.EXPECT().FindByCodes(gomock.Any(), []int{1}).Return(nil, infraErr)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := ports.UpdateStatusInput{
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

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusInDiagnosis, nil, nil, 0, "", ft(), ft())
	infraErr := errors.New("db timeout")

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	partRepo.EXPECT().FindByManufacturerCodes(gomock.Any(), []string{"FAB-001"}).Return(nil, infraErr)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := ports.UpdateStatusInput{
		ID:     "order-1",
		Status: entities.StatusAwaitingApproval,
		Parts: []ports.UpdateStatusPartInput{
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

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusInDiagnosis, nil, nil, 0, "", ft(), ft())
	infraErr := errors.New("db write error")

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(infraErr)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := ports.UpdateStatusInput{
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

	so := entities.ReconstituteServiceOrder("order-1", 0, "c1", "v1", entities.StatusReceived, nil, nil, 0, "", ft(), ft())
	infraErr := errors.New("db write error")

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	svcRepo := mocks.NewMockServiceRepository(ctrl)
	partRepo := mocks.NewMockPartRepository(ctrl)

	repo.EXPECT().FindByID(gomock.Any(), "order-1").Return(so, nil)
	repo.EXPECT().UpdateStatus(gomock.Any(), "order-1", entities.StatusInDiagnosis).Return(infraErr)

	uc := NewUpdateServiceOrderStatus(repo, svcRepo, partRepo)
	input := ports.UpdateStatusInput{
		ID:     "order-1",
		Status: entities.StatusInDiagnosis,
	}
	err := uc.Execute(context.Background(), input)
	if !errors.Is(err, infraErr) {
		t.Fatalf("erro = %v, esperava erro de infra propagado", err)
	}
}
