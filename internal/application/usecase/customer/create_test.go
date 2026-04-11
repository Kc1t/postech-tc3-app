package customeruc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
)

func fixedTime() time.Time { return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) }

// --- Mock do CustomerRepository ---

type mockCustomerRepo struct {
	findByDocumentFn func(ctx context.Context, doc string) (*entities.Customer, error)
	createFn         func(ctx context.Context, c *entities.Customer) error
}

func (m *mockCustomerRepo) Create(ctx context.Context, c *entities.Customer) error {
	return m.createFn(ctx, c)
}
func (m *mockCustomerRepo) FindByID(ctx context.Context, id string) (*entities.Customer, error) {
	return nil, domainerrors.ErrNotFound
}
func (m *mockCustomerRepo) FindByDocument(ctx context.Context, doc string) (*entities.Customer, error) {
	return m.findByDocumentFn(ctx, doc)
}
func (m *mockCustomerRepo) FindAll(ctx context.Context) ([]*entities.Customer, error) {
	return nil, nil
}
func (m *mockCustomerRepo) Update(ctx context.Context, c *entities.Customer) error { return nil }
func (m *mockCustomerRepo) Delete(ctx context.Context, id string) error            { return nil }

// --- Testes ---

func TestCreateCustomer_Sucesso(t *testing.T) {
	repo := &mockCustomerRepo{
		findByDocumentFn: func(_ context.Context, _ string) (*entities.Customer, error) {
			return nil, domainerrors.ErrNotFound
		},
		createFn: func(_ context.Context, _ *entities.Customer) error {
			return nil
		},
	}

	uc := NewCreateCustomer(repo)
	customer, err := entities.NewCustomer("Diego", "52998224725", "diego@email.com", "11999999999")
	if err != nil {
		t.Fatalf("erro ao criar customer: %v", err)
	}

	if err := uc.Execute(context.Background(), customer); err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
}

func TestCreateCustomer_DocumentoDuplicado(t *testing.T) {
	existing := entities.ReconstituteCustomer("id-1", "Outro", "52998224725", "outro@email.com", "", fixedTime(), fixedTime())
	repo := &mockCustomerRepo{
		findByDocumentFn: func(_ context.Context, _ string) (*entities.Customer, error) {
			return existing, nil
		},
		createFn: func(_ context.Context, _ *entities.Customer) error {
			t.Fatal("Create nao deveria ser chamado quando documento ja existe")
			return nil
		},
	}

	uc := NewCreateCustomer(repo)
	customer, _ := entities.NewCustomer("Diego", "52998224725", "diego@email.com", "")

	err := uc.Execute(context.Background(), customer)
	if err == nil {
		t.Fatal("esperava erro de duplicidade")
	}
	if !errors.Is(err, domainerrors.ErrAlreadyExists) {
		t.Errorf("erro = %v, esperava ErrAlreadyExists", err)
	}
}

func TestCreateCustomer_ErroNoPersistir(t *testing.T) {
	dbErr := errors.New("connection refused")
	repo := &mockCustomerRepo{
		findByDocumentFn: func(_ context.Context, _ string) (*entities.Customer, error) {
			return nil, domainerrors.ErrNotFound
		},
		createFn: func(_ context.Context, _ *entities.Customer) error {
			return dbErr
		},
	}

	uc := NewCreateCustomer(repo)
	customer, _ := entities.NewCustomer("Diego", "52998224725", "diego@email.com", "")

	err := uc.Execute(context.Background(), customer)
	if err == nil {
		t.Fatal("esperava erro do banco")
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("erro = %v, esperava %v", err, dbErr)
	}
}
