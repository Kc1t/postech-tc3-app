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

func TestListServiceOrdersByCPF_Sucesso(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	expected := []*entities.ServiceOrder{
		entities.NewServiceOrder("cust-1", "veh-1"),
	}

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil)
	soRepo.EXPECT().FindByCustomerID(gomock.Any(), "cust-1").Return(expected, nil)

	uc := NewListServiceOrdersByCPF(soRepo, custRepo)
	got, err := uc.Execute(context.Background(), "52998224725")
	if err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, esperava 1", len(got))
	}
}

func TestListServiceOrdersByCPF_ListaVazia(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil)
	soRepo.EXPECT().FindByCustomerID(gomock.Any(), "cust-1").Return([]*entities.ServiceOrder{}, nil)

	uc := NewListServiceOrdersByCPF(soRepo, custRepo)
	got, err := uc.Execute(context.Background(), "52998224725")
	if err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("len = %d, esperava 0", len(got))
	}
}

func TestListServiceOrdersByCPF_ClienteNaoEncontrado(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "00000000000").Return(nil, domainerrors.ErrNotFound)

	uc := NewListServiceOrdersByCPF(soRepo, custRepo)
	_, err := uc.Execute(context.Background(), "00000000000")
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("erro = %v, esperava ErrNotFound", err)
	}
}

func TestListServiceOrdersByCPF_ErroInfraCliente(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	infraErr := errors.New("timeout")

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(nil, infraErr)

	uc := NewListServiceOrdersByCPF(soRepo, custRepo)
	_, err := uc.Execute(context.Background(), "52998224725")
	if !errors.Is(err, infraErr) {
		t.Fatalf("erro = %v, esperava erro de infra propagado", err)
	}
}

func TestListServiceOrdersByCPF_ErroInfraRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	infraErr := errors.New("db error")

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil)
	soRepo.EXPECT().FindByCustomerID(gomock.Any(), "cust-1").Return(nil, infraErr)

	uc := NewListServiceOrdersByCPF(soRepo, custRepo)
	_, err := uc.Execute(context.Background(), "52998224725")
	if !errors.Is(err, infraErr) {
		t.Fatalf("erro = %v, esperava erro de infra propagado", err)
	}
}
