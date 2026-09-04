package entities

import (
	"testing"
	"time"
)

func TestNewPart(t *testing.T) {
	before := time.Now()
	p := NewPart("FAB-001", "Filtro de óleo", "Filtro original", "unidade", 49.90, 10)
	after := time.Now()

	if p.Name() != "Filtro de óleo" {
		t.Errorf("expected name %q, got %q", "Filtro de óleo", p.Name())
	}
	if p.Description() != "Filtro original" {
		t.Errorf("expected description %q, got %q", "Filtro original", p.Description())
	}
	if p.Unit() != "unidade" {
		t.Errorf("expected unit %q, got %q", "unidade", p.Unit())
	}
	if p.Price() != 49.90 {
		t.Errorf("expected price %v, got %v", 49.90, p.Price())
	}
	if p.Stock() != 10 {
		t.Errorf("expected stock %d, got %d", 10, p.Stock())
	}
	if p.ID() != "" {
		t.Errorf("expected empty ID, got %q", p.ID())
	}
	if p.CreatedAt().Before(before) || p.CreatedAt().After(after) {
		t.Error("CreatedAt out of expected range")
	}
}

func TestReconstitutePart(t *testing.T) {
	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	p := ReconstitutePart("pid-1", "FAB-001", "Filtro", "Desc", "unidade", 30.0, 5, createdAt, updatedAt)

	if p.ID() != "pid-1" {
		t.Errorf("expected ID %q, got %q", "pid-1", p.ID())
	}
	if p.Stock() != 5 {
		t.Errorf("expected stock %d, got %d", 5, p.Stock())
	}
	if !p.CreatedAt().Equal(createdAt) {
		t.Errorf("expected createdAt %v, got %v", createdAt, p.CreatedAt())
	}
}

func TestPart_SetID(t *testing.T) {
	p := NewPart("FAB-002", "Filtro", "Desc", "unidade", 10.0, 5)
	p.SetID("pid-1")
	if p.ID() != "pid-1" {
		t.Errorf("expected ID %q, got %q", "pid-1", p.ID())
	}
}

func TestPart_SetName_TouchesUpdatedAt(t *testing.T) {
	p := NewPart("FAB-002", "Filtro", "Desc", "unidade", 10.0, 5)
	before := p.UpdatedAt()
	time.Sleep(time.Millisecond)
	p.SetName("Vela")

	if p.Name() != "Vela" {
		t.Errorf("expected name %q, got %q", "Vela", p.Name())
	}
	if !p.UpdatedAt().After(before) {
		t.Error("expected UpdatedAt to be updated after SetName")
	}
}

func TestPart_SetDescription_TouchesUpdatedAt(t *testing.T) {
	p := NewPart("FAB-002", "Filtro", "Desc", "unidade", 10.0, 5)
	before := p.UpdatedAt()
	time.Sleep(time.Millisecond)
	p.SetDescription("Nova desc")

	if p.Description() != "Nova desc" {
		t.Errorf("expected description %q, got %q", "Nova desc", p.Description())
	}
	if !p.UpdatedAt().After(before) {
		t.Error("expected UpdatedAt to be updated after SetDescription")
	}
}

func TestPart_SetUnit_TouchesUpdatedAt(t *testing.T) {
	p := NewPart("FAB-002", "Filtro", "Desc", "unidade", 10.0, 5)
	before := p.UpdatedAt()
	time.Sleep(time.Millisecond)
	p.SetUnit("litro")

	if p.Unit() != "litro" {
		t.Errorf("expected unit %q, got %q", "litro", p.Unit())
	}
	if !p.UpdatedAt().After(before) {
		t.Error("expected UpdatedAt to be updated after SetUnit")
	}
}

func TestPart_SetPrice_TouchesUpdatedAt(t *testing.T) {
	p := NewPart("FAB-002", "Filtro", "Desc", "unidade", 10.0, 5)
	before := p.UpdatedAt()
	time.Sleep(time.Millisecond)
	p.SetPrice(99.99)

	if p.Price() != 99.99 {
		t.Errorf("expected price %v, got %v", 99.99, p.Price())
	}
	if !p.UpdatedAt().After(before) {
		t.Error("expected UpdatedAt to be updated after SetPrice")
	}
}

func TestPart_SetStock_TouchesUpdatedAt(t *testing.T) {
	p := NewPart("FAB-002", "Filtro", "Desc", "unidade", 10.0, 5)
	before := p.UpdatedAt()
	time.Sleep(time.Millisecond)
	p.SetStock(20)

	if p.Stock() != 20 {
		t.Errorf("expected stock %d, got %d", 20, p.Stock())
	}
	if !p.UpdatedAt().After(before) {
		t.Error("expected UpdatedAt to be updated after SetStock")
	}
}
