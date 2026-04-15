package customeruc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestCreateCustomer_Sucesso(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockCustomerRepository(ctrl)
	repo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(nil, domainerrors.ErrNotFound)
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	existing := entities.ReconstituteCustomer("id-1", "Outro", "52998224725", "outro@email.com", "", time.Now(), time.Now())

	repo := mocks.NewMockCustomerRepository(ctrl)
	repo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(existing, nil)
	// Create NAO deve ser chamado

	uc := NewCreateCustomer(repo)
	customer, _ := entities.NewCustomer("Diego", "52998224725", "diego@email.com", "")

	err := uc.Execute(context.Background(), customer)
	if !errors.Is(err, domainerrors.ErrAlreadyExists) {
		t.Fatalf("erro = %v, esperava ErrAlreadyExists", err)
	}
}

func TestCreateCustomer_ErroNoPersistir(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbErr := errors.New("connection refused")

	repo := mocks.NewMockCustomerRepository(ctrl)
	repo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(nil, domainerrors.ErrNotFound)
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr)

	uc := NewCreateCustomer(repo)
	customer, _ := entities.NewCustomer("Diego", "52998224725", "diego@email.com", "")

	err := uc.Execute(context.Background(), customer)
	if !errors.Is(err, dbErr) {
		t.Fatalf("erro = %v, esperava %v", err, dbErr)
	}
}

func TestCreateCustomer_ErroInfraNoFindByDocument(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	infraErr := errors.New("timeout")

	repo := mocks.NewMockCustomerRepository(ctrl)
	repo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(nil, infraErr)
	// Create NAO deve ser chamado

	uc := NewCreateCustomer(repo)
	customer, _ := entities.NewCustomer("Diego", "52998224725", "diego@email.com", "")

	err := uc.Execute(context.Background(), customer)
	if !errors.Is(err, infraErr) {
		t.Fatalf("erro = %v, esperava %v (erro de infra propagado)", err, infraErr)
	}
}
