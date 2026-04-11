package serviceorderuc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
)

// --- Mock com FindByID configuravel ---

type mockSORepoForStatus struct {
	findByIDFn   func(ctx context.Context, id string) (*entities.ServiceOrder, error)
	updateCalled bool
	updateErr    error
}

func (m *mockSORepoForStatus) Create(ctx context.Context, so *entities.ServiceOrder) error {
	return nil
}
func (m *mockSORepoForStatus) FindByID(ctx context.Context, id string) (*entities.ServiceOrder, error) {
	return m.findByIDFn(ctx, id)
}
func (m *mockSORepoForStatus) FindAll(ctx context.Context) ([]*entities.ServiceOrder, error) {
	return nil, nil
}
func (m *mockSORepoForStatus) FindByCustomerID(ctx context.Context, id string) ([]*entities.ServiceOrder, error) {
	return nil, nil
}
func (m *mockSORepoForStatus) UpdateStatus(ctx context.Context, id string, s entities.OrderStatus) error {
	m.updateCalled = true
	return m.updateErr
}
func (m *mockSORepoForStatus) Update(ctx context.Context, so *entities.ServiceOrder) error {
	return nil
}
func (m *mockSORepoForStatus) Delete(ctx context.Context, id string) error { return nil }

// --- Testes ---

func TestUpdateStatus_TransicaoValida(t *testing.T) {
	now := time.Now()
	so := entities.ReconstituteServiceOrder("order-1", "c1", "v1", entities.StatusReceived, nil, nil, 0, "", now, now)

	repo := &mockSORepoForStatus{
		findByIDFn: func(_ context.Context, _ string) (*entities.ServiceOrder, error) {
			return so, nil
		},
	}

	uc := NewUpdateServiceOrderStatus(repo)
	err := uc.Execute(context.Background(), "order-1", entities.StatusInDiagnosis)
	if err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
	if !repo.updateCalled {
		t.Error("UpdateStatus do repo deveria ter sido chamado")
	}
}

func TestUpdateStatus_TransicaoInvalida(t *testing.T) {
	now := time.Now()
	so := entities.ReconstituteServiceOrder("order-1", "c1", "v1", entities.StatusReceived, nil, nil, 0, "", now, now)

	repo := &mockSORepoForStatus{
		findByIDFn: func(_ context.Context, _ string) (*entities.ServiceOrder, error) {
			return so, nil
		},
	}

	uc := NewUpdateServiceOrderStatus(repo)
	err := uc.Execute(context.Background(), "order-1", entities.StatusFinished)
	if err == nil {
		t.Fatal("esperava erro de transicao invalida")
	}
	if !errors.Is(err, domainerrors.ErrInvalidStatus) {
		t.Errorf("erro = %v, esperava ErrInvalidStatus", err)
	}
	if repo.updateCalled {
		t.Error("UpdateStatus do repo NAO deveria ter sido chamado em transicao invalida")
	}
}

func TestUpdateStatus_OSNaoEncontrada(t *testing.T) {
	repo := &mockSORepoForStatus{
		findByIDFn: func(_ context.Context, _ string) (*entities.ServiceOrder, error) {
			return nil, domainerrors.ErrNotFound
		},
	}

	uc := NewUpdateServiceOrderStatus(repo)
	err := uc.Execute(context.Background(), "inexistente", entities.StatusInDiagnosis)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Errorf("erro = %v, esperava ErrNotFound", err)
	}
}

func TestUpdateStatus_StatusDesconhecido(t *testing.T) {
	now := time.Now()
	so := entities.ReconstituteServiceOrder("order-1", "c1", "v1", entities.StatusReceived, nil, nil, 0, "", now, now)

	repo := &mockSORepoForStatus{
		findByIDFn: func(_ context.Context, _ string) (*entities.ServiceOrder, error) {
			return so, nil
		},
	}

	uc := NewUpdateServiceOrderStatus(repo)
	err := uc.Execute(context.Background(), "order-1", entities.OrderStatus("invalido"))
	if err == nil {
		t.Fatal("esperava erro para status desconhecido")
	}
	if repo.updateCalled {
		t.Error("UpdateStatus do repo NAO deveria ter sido chamado")
	}
}
