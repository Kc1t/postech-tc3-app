package postgresql

import (
	"context"
	"testing"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/entities"
)

func TestServiceRepository_Create(t *testing.T) {
	repo := NewServiceRepository(testDB)
	s := entities.NewService(1, "Troca de óleo", "Troca de óleo sintético", 150.00, 60)

	if err := repo.Create(context.Background(), s); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if s.ID() == "" {
		t.Fatal("expected ID to be set after create")
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Service{}, "id = ?", s.ID()) })
}

func TestServiceRepository_FindByID(t *testing.T) {
	repo := NewServiceRepository(testDB)
	s := entities.NewService(2, "Alinhamento", "Alinhamento de direção", 80.00, 30)
	if err := repo.Create(context.Background(), s); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Service{}, "id = ?", s.ID()) })

	found, err := repo.FindByID(context.Background(), s.ID())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.Name() != "Alinhamento" {
		t.Errorf("expected name Alinhamento, got %s", found.Name())
	}
}

func TestServiceRepository_FindByID_NotFound(t *testing.T) {
	repo := NewServiceRepository(testDB)

	_, err := repo.FindByID(context.Background(), "00000000-0000-0000-0000-000000000000")
	if !isDomainNotFound(err) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestServiceRepository_FindAll(t *testing.T) {
	repo := NewServiceRepository(testDB)
	s := entities.NewService(3, "Balanceamento", "Balanceamento de rodas", 60.00, 45)
	if err := repo.Create(context.Background(), s); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Service{}, "id = ?", s.ID()) })

	all, err := repo.FindAll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(all) == 0 {
		t.Fatal("expected at least one service")
	}
}

func TestServiceRepository_Update(t *testing.T) {
	repo := NewServiceRepository(testDB)
	s := entities.NewService(4, "Revisão 10k", "Revisão de 10.000 km", 350.00, 120)
	if err := repo.Create(context.Background(), s); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Service{}, "id = ?", s.ID()) })

	s.SetPrice(380.00)
	if err := repo.Update(context.Background(), s); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, _ := repo.FindByID(context.Background(), s.ID())
	if found.Price() != 380.00 {
		t.Errorf("expected price 380.00, got %f", found.Price())
	}
}

func TestServiceRepository_Delete(t *testing.T) {
	repo := NewServiceRepository(testDB)
	s := entities.NewService(5, "Higienização", "Higienização interna", 200.00, 90)
	if err := repo.Create(context.Background(), s); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if err := repo.Delete(context.Background(), s.ID()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err := repo.FindByID(context.Background(), s.ID())
	if !isDomainNotFound(err) {
		t.Fatal("expected record to be deleted")
	}
}

func TestServiceRepository_FindByIDs(t *testing.T) {
	repo := NewServiceRepository(testDB)
	s1 := entities.NewService(101, "Diagnostico eletrico", "Scanner OBD2", 180.00, 45)
	s2 := entities.NewService(102, "Troca de correia", "Correia dentada", 320.00, 90)
	if err := repo.Create(context.Background(), s1); err != nil {
		t.Fatalf("setup s1 failed: %v", err)
	}
	if err := repo.Create(context.Background(), s2); err != nil {
		t.Fatalf("setup s2 failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Service{}, "id IN ?", []string{s1.ID(), s2.ID()}) })

	found, err := repo.FindByIDs(context.Background(), []string{s1.ID(), s2.ID()})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(found) != 2 {
		t.Errorf("esperava 2 servicos, obteve %d", len(found))
	}
}

func TestServiceRepository_FindByCodes(t *testing.T) {
	repo := NewServiceRepository(testDB)
	s1 := entities.NewService(201, "Lavagem tecnica", "Lavagem detalhada", 90.00, 60)
	s2 := entities.NewService(202, "Polimento", "Polimento completo", 250.00, 180)
	if err := repo.Create(context.Background(), s1); err != nil {
		t.Fatalf("setup s1 failed: %v", err)
	}
	if err := repo.Create(context.Background(), s2); err != nil {
		t.Fatalf("setup s2 failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Service{}, "id IN ?", []string{s1.ID(), s2.ID()}) })

	found, err := repo.FindByCodes(context.Background(), []int{201, 202, 9999})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Retorna so os servicos que existem (201 e 202); 9999 e ignorado.
	if len(found) != 2 {
		t.Errorf("esperava 2 servicos, obteve %d", len(found))
	}
}
