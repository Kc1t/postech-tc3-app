package entities

import (
	"testing"
	"time"
)

func TestNewService(t *testing.T) {
	before := time.Now()
	s := NewService("Troca de óleo", "Troca completa", 150.0, 60)
	after := time.Now()

	if s.Name() != "Troca de óleo" {
		t.Errorf("expected name %q, got %q", "Troca de óleo", s.Name())
	}
	if s.Description() != "Troca completa" {
		t.Errorf("expected description %q, got %q", "Troca completa", s.Description())
	}
	if s.Price() != 150.0 {
		t.Errorf("expected price %v, got %v", 150.0, s.Price())
	}
	if s.DurationMin() != 60 {
		t.Errorf("expected durationMin %d, got %d", 60, s.DurationMin())
	}
	if s.ID() != "" {
		t.Errorf("expected empty ID, got %q", s.ID())
	}
	if s.CreatedAt().Before(before) || s.CreatedAt().After(after) {
		t.Error("CreatedAt out of expected range")
	}
}

func TestReconstituteService(t *testing.T) {
	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	s := ReconstituteService("sid-1", "Alinhamento", "Desc", 200.0, 90, createdAt, updatedAt)

	if s.ID() != "sid-1" {
		t.Errorf("expected ID %q, got %q", "sid-1", s.ID())
	}
	if s.DurationMin() != 90 {
		t.Errorf("expected durationMin %d, got %d", 90, s.DurationMin())
	}
	if !s.UpdatedAt().Equal(updatedAt) {
		t.Errorf("expected updatedAt %v, got %v", updatedAt, s.UpdatedAt())
	}
}

func TestService_SetID(t *testing.T) {
	s := NewService("Troca de óleo", "Desc", 100.0, 30)
	s.SetID("sid-1")
	if s.ID() != "sid-1" {
		t.Errorf("expected ID %q, got %q", "sid-1", s.ID())
	}
}

func TestService_SetName_TouchesUpdatedAt(t *testing.T) {
	s := NewService("Troca de óleo", "Desc", 100.0, 30)
	before := s.UpdatedAt()
	time.Sleep(time.Millisecond)
	s.SetName("Alinhamento")

	if s.Name() != "Alinhamento" {
		t.Errorf("expected name %q, got %q", "Alinhamento", s.Name())
	}
	if !s.UpdatedAt().After(before) {
		t.Error("expected UpdatedAt to be updated after SetName")
	}
}

func TestService_SetDescription_TouchesUpdatedAt(t *testing.T) {
	s := NewService("Troca de óleo", "Desc", 100.0, 30)
	before := s.UpdatedAt()
	time.Sleep(time.Millisecond)
	s.SetDescription("Nova descrição")

	if s.Description() != "Nova descrição" {
		t.Errorf("expected description %q, got %q", "Nova descrição", s.Description())
	}
	if !s.UpdatedAt().After(before) {
		t.Error("expected UpdatedAt to be updated after SetDescription")
	}
}

func TestService_SetPrice_TouchesUpdatedAt(t *testing.T) {
	s := NewService("Troca de óleo", "Desc", 100.0, 30)
	before := s.UpdatedAt()
	time.Sleep(time.Millisecond)
	s.SetPrice(250.0)

	if s.Price() != 250.0 {
		t.Errorf("expected price %v, got %v", 250.0, s.Price())
	}
	if !s.UpdatedAt().After(before) {
		t.Error("expected UpdatedAt to be updated after SetPrice")
	}
}

func TestService_SetDurationMin_TouchesUpdatedAt(t *testing.T) {
	s := NewService("Troca de óleo", "Desc", 100.0, 30)
	before := s.UpdatedAt()
	time.Sleep(time.Millisecond)
	s.SetDurationMin(120)

	if s.DurationMin() != 120 {
		t.Errorf("expected durationMin %d, got %d", 120, s.DurationMin())
	}
	if !s.UpdatedAt().After(before) {
		t.Error("expected UpdatedAt to be updated after SetDurationMin")
	}
}
