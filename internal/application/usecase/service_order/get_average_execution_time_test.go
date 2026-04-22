package serviceorderuc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fiap/postech-tc1/internal/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestGetAverageExecutionTime_Sucesso(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := 90 * time.Minute
	repo := mocks.NewMockServiceOrderRepository(ctrl)
	repo.EXPECT().AverageExecutionTime(gomock.Any()).Return(expected, nil)

	uc := NewGetAverageExecutionTime(repo)
	got, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
	if got != expected {
		t.Errorf("media = %v, esperava %v", got, expected)
	}
}

func TestGetAverageExecutionTime_SemAmostras(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockServiceOrderRepository(ctrl)
	repo.EXPECT().AverageExecutionTime(gomock.Any()).Return(time.Duration(0), nil)

	uc := NewGetAverageExecutionTime(repo)
	got, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
	if got != 0 {
		t.Errorf("media = %v, esperava 0 (sem amostras)", got)
	}
}

func TestGetAverageExecutionTime_ErroInfra(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	infraErr := errors.New("db unavailable")
	repo := mocks.NewMockServiceOrderRepository(ctrl)
	repo.EXPECT().AverageExecutionTime(gomock.Any()).Return(time.Duration(0), infraErr)

	uc := NewGetAverageExecutionTime(repo)
	_, err := uc.Execute(context.Background())
	if !errors.Is(err, infraErr) {
		t.Fatalf("erro = %v, esperava erro de infra propagado", err)
	}
}
