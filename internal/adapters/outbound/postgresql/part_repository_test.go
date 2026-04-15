package postgresql

import (
	"context"
	"testing"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/entities"
)

func TestPartRepository_Create(t *testing.T) {
	repo := NewPartRepository(testDB)
	p := entities.NewPart("FAB-001", "Filtro de óleo", "Filtro para motor", "unidade", 49.90, 10)

	if err := repo.Create(context.Background(), p); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.ID() == "" {
		t.Fatal("expected ID to be set after create")
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Part{}, "id = ?", p.ID()) })
}

func TestPartRepository_FindByID(t *testing.T) {
	repo := NewPartRepository(testDB)
	p := entities.NewPart("FAB-002", "Vela de ignição", "Vela NGK", "unidade", 25.00, 8)
	if err := repo.Create(context.Background(), p); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Part{}, "id = ?", p.ID()) })

	found, err := repo.FindByID(context.Background(), p.ID())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.Name() != "Vela de ignição" {
		t.Errorf("expected name Vela de ignição, got %s", found.Name())
	}
}

func TestPartRepository_FindByID_NotFound(t *testing.T) {
	repo := NewPartRepository(testDB)

	_, err := repo.FindByID(context.Background(), "00000000-0000-0000-0000-000000000000")
	if !isDomainNotFound(err) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPartRepository_FindAll(t *testing.T) {
	repo := NewPartRepository(testDB)
	p := entities.NewPart("FAB-003", "Correia dentada", "Correia Gates", "unidade", 120.00, 3)
	if err := repo.Create(context.Background(), p); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Part{}, "id = ?", p.ID()) })

	all, err := repo.FindAll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(all) == 0 {
		t.Fatal("expected at least one part")
	}
}

func TestPartRepository_Update(t *testing.T) {
	repo := NewPartRepository(testDB)
	p := entities.NewPart("FAB-004", "Amortecedor", "Amortecedor Cofap", "unidade", 250.00, 4)
	if err := repo.Create(context.Background(), p); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Part{}, "id = ?", p.ID()) })

	p.SetPrice(280.00)
	if err := repo.Update(context.Background(), p); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, _ := repo.FindByID(context.Background(), p.ID())
	if found.Price() != 280.00 {
		t.Errorf("expected price 280.00, got %f", found.Price())
	}
}

func TestPartRepository_Delete(t *testing.T) {
	repo := NewPartRepository(testDB)
	p := entities.NewPart("FAB-005", "Palheta do limpador", "Bosch 21", "unidade", 35.00, 6)
	if err := repo.Create(context.Background(), p); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if err := repo.Delete(context.Background(), p.ID()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err := repo.FindByID(context.Background(), p.ID())
	if !isDomainNotFound(err) {
		t.Fatal("expected record to be deleted")
	}
}

func TestPartRepository_UpdateStock(t *testing.T) {
	repo := NewPartRepository(testDB)
	p := entities.NewPart("FAB-006", "Pastilha de freio", "Fremax", "jogo", 180.00, 5)
	if err := repo.Create(context.Background(), p); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Part{}, "id = ?", p.ID()) })

	if err := repo.UpdateStock(context.Background(), p.ID(), 3); err != nil {
		t.Fatalf("expected no error on increment, got %v", err)
	}

	found, _ := repo.FindByID(context.Background(), p.ID())
	if found.Stock() != 8 {
		t.Errorf("expected stock 8, got %d", found.Stock())
	}

	if err := repo.UpdateStock(context.Background(), p.ID(), -2); err != nil {
		t.Fatalf("expected no error on decrement, got %v", err)
	}

	found, _ = repo.FindByID(context.Background(), p.ID())
	if found.Stock() != 6 {
		t.Errorf("expected stock 6, got %d", found.Stock())
	}
}
