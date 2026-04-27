package postgresql

import (
	"context"
	"errors"
	"testing"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
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

func TestPartRepository_UpdateStock_ZeroesStock(t *testing.T) {
	repo := NewPartRepository(testDB)
	p := entities.NewPart("FAB-007", "Filtro de ar", "Tecfil", "unidade", 40.00, 4)
	if err := repo.Create(context.Background(), p); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Part{}, "id = ?", p.ID()) })

	if err := repo.UpdateStock(context.Background(), p.ID(), -4); err != nil {
		t.Fatalf("expected no error zeroing stock, got %v", err)
	}

	found, _ := repo.FindByID(context.Background(), p.ID())
	if found.Stock() != 0 {
		t.Errorf("expected stock 0, got %d", found.Stock())
	}
}

func TestPartRepository_UpdateStock_InsufficientStock(t *testing.T) {
	repo := NewPartRepository(testDB)
	p := entities.NewPart("FAB-008", "Fluido de freio", "Bosch DOT4", "litro", 30.00, 2)
	if err := repo.Create(context.Background(), p); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Part{}, "id = ?", p.ID()) })

	err := repo.UpdateStock(context.Background(), p.ID(), -5)
	if !errors.Is(err, domainerrors.ErrInsufficientStock) {
		t.Fatalf("expected ErrInsufficientStock, got %v", err)
	}

	found, _ := repo.FindByID(context.Background(), p.ID())
	if found.Stock() != 2 {
		t.Errorf("expected stock unchanged (2), got %d", found.Stock())
	}
}

func TestPartRepository_UpdateStock_NotFound(t *testing.T) {
	repo := NewPartRepository(testDB)
	err := repo.UpdateStock(context.Background(), "00000000-0000-0000-0000-000000000000", 1)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPartRepository_FindByIDs(t *testing.T) {
	repo := NewPartRepository(testDB)
	p1 := entities.NewPart("FAB-100", "Filtro combustivel", "Filtro", "unidade", 35.00, 20)
	p2 := entities.NewPart("FAB-101", "Oleo sintetico 5W30", "1L", "litro", 65.00, 40)
	if err := repo.Create(context.Background(), p1); err != nil {
		t.Fatalf("setup p1 failed: %v", err)
	}
	if err := repo.Create(context.Background(), p2); err != nil {
		t.Fatalf("setup p2 failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Part{}, "id IN ?", []string{p1.ID(), p2.ID()}) })

	found, err := repo.FindByIDs(context.Background(), []string{p1.ID(), p2.ID()})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(found) != 2 {
		t.Errorf("esperava 2 pecas, obteve %d", len(found))
	}
}

func TestPartRepository_FindByManufacturerCodes(t *testing.T) {
	repo := NewPartRepository(testDB)
	p1 := entities.NewPart("MFG-AAA", "Velas de ignicao", "NGK", "jogo", 120.00, 15)
	p2 := entities.NewPart("MFG-BBB", "Pastilhas dianteiras", "Fremax", "jogo", 210.00, 8)
	if err := repo.Create(context.Background(), p1); err != nil {
		t.Fatalf("setup p1 failed: %v", err)
	}
	if err := repo.Create(context.Background(), p2); err != nil {
		t.Fatalf("setup p2 failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Part{}, "id IN ?", []string{p1.ID(), p2.ID()}) })

	found, err := repo.FindByManufacturerCodes(context.Background(), []string{"MFG-AAA", "MFG-BBB", "MFG-INEXISTENTE"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(found) != 2 {
		t.Errorf("esperava 2 pecas, obteve %d", len(found))
	}
}
