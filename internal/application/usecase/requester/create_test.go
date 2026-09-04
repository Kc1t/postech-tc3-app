package requesteruc

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

func TestCreateRequester_Sucesso(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRequesterRepository(ctrl)
	repo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(nil, domainerrors.ErrNotFound)
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	uc := NewCreateRequester(repo)
	requester, err := entities.NewRequester("Diego", "52998224725", "diego@email.com", "11999999999")
	if err != nil {
		t.Fatalf("erro ao criar requester: %v", err)
	}

	if err := uc.Execute(context.Background(), requester); err != nil {
		t.Fatalf("esperava sucesso, obteve: %v", err)
	}
}

func TestCreateRequester_DocumentoDuplicado(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	existing := entities.ReconstituteRequester("id-1", "Outro", "52998224725", "outro@email.com", "", time.Now(), time.Now())

	repo := mocks.NewMockRequesterRepository(ctrl)
	repo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(existing, nil)

	uc := NewCreateRequester(repo)
	requester, _ := entities.NewRequester("Diego", "52998224725", "diego@email.com", "")

	err := uc.Execute(context.Background(), requester)
	if !errors.Is(err, domainerrors.ErrAlreadyExists) {
		t.Fatalf("erro = %v, esperava ErrAlreadyExists", err)
	}
}

func TestCreateRequester_ErroNoPersistir(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbErr := errors.New("connection refused")

	repo := mocks.NewMockRequesterRepository(ctrl)
	repo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(nil, domainerrors.ErrNotFound)
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr)

	uc := NewCreateRequester(repo)
	requester, _ := entities.NewRequester("Diego", "52998224725", "diego@email.com", "")

	err := uc.Execute(context.Background(), requester)
	if !errors.Is(err, dbErr) {
		t.Fatalf("erro = %v, esperava %v", err, dbErr)
	}
}

func TestCreateRequester_ErroInfraNoFindByDocument(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	infraErr := errors.New("timeout")

	repo := mocks.NewMockRequesterRepository(ctrl)
	repo.EXPECT().FindByDocument(gomock.Any(), "52998224725").Return(nil, infraErr)

	uc := NewCreateRequester(repo)
	requester, _ := entities.NewRequester("Diego", "52998224725", "diego@email.com", "")

	err := uc.Execute(context.Background(), requester)
	if !errors.Is(err, infraErr) {
		t.Fatalf("erro = %v, esperava %v (erro de infra propagado)", err, infraErr)
	}
}
