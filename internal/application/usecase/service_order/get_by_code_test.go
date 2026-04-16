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

func TestGetServiceOrderByCode_Sucesso(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	expected := entities.ReconstituteServiceOrder("order-1", 100, "cust-1", "veh-1",
		entities.StatusReceived, nil, nil, 0, "", ft(), ft())

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil)
	soRepo.EXPECT().FindByCode(gomock.Any(), 100).Return(expected, nil)

	uc := NewGetServiceOrderByCode(soRepo, custRepo)
	got, err := uc.Execute(context.Background(), 100, "52998224725")
	if err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
	if got.ID() != "order-1" {
		t.Errorf("ID = %q, esperava %q", got.ID(), "order-1")
	}
}

func TestGetServiceOrderByCode_ClienteNaoEncontrado(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "00000000000").Return(nil, domainerrors.ErrNotFound)

	uc := NewGetServiceOrderByCode(soRepo, custRepo)
	_, err := uc.Execute(context.Background(), 100, "00000000000")
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("erro = %v, esperava ErrNotFound", err)
	}
}

func TestGetServiceOrderByCode_OSNaoEncontrada(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil)
	soRepo.EXPECT().FindByCode(gomock.Any(), 999).Return(nil, domainerrors.ErrNotFound)

	uc := NewGetServiceOrderByCode(soRepo, custRepo)
	_, err := uc.Execute(context.Background(), 999, "52998224725")
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("erro = %v, esperava ErrNotFound", err)
	}
}

func TestGetServiceOrderByCode_OSDeOutroCliente(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customer := entities.ReconstituteCustomer("cust-1", "Diego", "52998224725", "d@e.com", "", ft(), ft())
	so := entities.ReconstituteServiceOrder("order-1", 100, "outro-cliente", "veh-1",
		entities.StatusReceived, nil, nil, 0, "", ft(), ft())

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(customer, nil)
	soRepo.EXPECT().FindByCode(gomock.Any(), 100).Return(so, nil)

	uc := NewGetServiceOrderByCode(soRepo, custRepo)
	_, err := uc.Execute(context.Background(), 100, "52998224725")
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("erro = %v, esperava ErrNotFound (ownership)", err)
	}
}

func TestGetServiceOrderByCode_ErroInfraCliente(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	infraErr := errors.New("timeout")

	soRepo := mocks.NewMockServiceOrderRepository(ctrl)
	custRepo := mocks.NewMockCustomerRepository(ctrl)

	custRepo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(nil, infraErr)

	uc := NewGetServiceOrderByCode(soRepo, custRepo)
	_, err := uc.Execute(context.Background(), 100, "52998224725")
	if !errors.Is(err, infraErr) {
		t.Fatalf("erro = %v, esperava erro de infra propagado", err)
	}
}
