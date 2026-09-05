package entities

import (
	"testing"
	"time"
)

func TestNewRequester(t *testing.T) {
	before := time.Now()
	c, err := NewRequester("João Silva", "52998224725", "joao@email.com", "11999990000")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	after := time.Now()

	if c.Name() != "João Silva" {
		t.Errorf("expected name %q, got %q", "João Silva", c.Name())
	}
	if c.Document() != "52998224725" {
		t.Errorf("expected document %q, got %q", "52998224725", c.Document())
	}
	if c.Email() != "joao@email.com" {
		t.Errorf("expected email %q, got %q", "joao@email.com", c.Email())
	}
	if c.Phone() != "11999990000" {
		t.Errorf("expected phone %q, got %q", "11999990000", c.Phone())
	}
	if c.ID() != "" {
		t.Errorf("expected empty ID, got %q", c.ID())
	}
	if c.CreatedAt().Before(before) || c.CreatedAt().After(after) {
		t.Error("CreatedAt out of expected range")
	}
	if c.UpdatedAt().Before(before) || c.UpdatedAt().After(after) {
		t.Error("UpdatedAt out of expected range")
	}
}

func TestNewRequester_InvalidDocument(t *testing.T) {
	_, err := NewRequester("João", "123", "j@j.com", "11999")
	if err == nil {
		t.Fatal("expected error for invalid document, got nil")
	}
}

func TestReconstituteRequester(t *testing.T) {
	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	c := ReconstituteRequester("id-1", "Maria", "98765432100", "maria@email.com", "11888880000", createdAt, updatedAt)

	if c.ID() != "id-1" {
		t.Errorf("expected ID %q, got %q", "id-1", c.ID())
	}
	if c.Name() != "Maria" {
		t.Errorf("expected name %q, got %q", "Maria", c.Name())
	}
	if !c.CreatedAt().Equal(createdAt) {
		t.Errorf("expected createdAt %v, got %v", createdAt, c.CreatedAt())
	}
	if !c.UpdatedAt().Equal(updatedAt) {
		t.Errorf("expected updatedAt %v, got %v", updatedAt, c.UpdatedAt())
	}
}

func TestRequester_SetID(t *testing.T) {
	c, err := NewRequester("João", "52998224725", "j@j.com", "11999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c.SetID("new-id")
	if c.ID() != "new-id" {
		t.Errorf("expected ID %q, got %q", "new-id", c.ID())
	}
}

func TestRequester_SetName_TouchesUpdatedAt(t *testing.T) {
	c, err := NewRequester("João", "52998224725", "j@j.com", "11999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	before := c.UpdatedAt()
	time.Sleep(time.Millisecond)
	c.SetName("Pedro")

	if c.Name() != "Pedro" {
		t.Errorf("expected name %q, got %q", "Pedro", c.Name())
	}
	if !c.UpdatedAt().After(before) {
		t.Error("expected UpdatedAt to be updated after SetName")
	}
}

func TestRequester_SetEmail_TouchesUpdatedAt(t *testing.T) {
	c, err := NewRequester("João", "52998224725", "j@j.com", "11999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	before := c.UpdatedAt()
	time.Sleep(time.Millisecond)
	c.SetEmail("novo@email.com")

	if c.Email() != "novo@email.com" {
		t.Errorf("expected email %q, got %q", "novo@email.com", c.Email())
	}
	if !c.UpdatedAt().After(before) {
		t.Error("expected UpdatedAt to be updated after SetEmail")
	}
}

func TestRequester_SetPhone_TouchesUpdatedAt(t *testing.T) {
	c, err := NewRequester("João", "52998224725", "j@j.com", "11999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	before := c.UpdatedAt()
	time.Sleep(time.Millisecond)
	c.SetPhone("11888880000")

	if c.Phone() != "11888880000" {
		t.Errorf("expected phone %q, got %q", "11888880000", c.Phone())
	}
	if !c.UpdatedAt().After(before) {
		t.Error("expected UpdatedAt to be updated after SetPhone")
	}
}

func TestNewRequester_StartsActive(t *testing.T) {
	c, err := NewRequester("João", "52998224725", "j@j.com", "11999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.Status() != RequesterStatusActive {
		t.Errorf("expected status %q, got %q", RequesterStatusActive, c.Status())
	}
	if !c.IsActive() {
		t.Error("expected new requester to be active")
	}
}

func TestReconstituteRequester_StartsActive(t *testing.T) {
	now := time.Now()
	c := ReconstituteRequester("id-1", "Maria", "98765432100", "maria@email.com", "11888880000", now, now)

	if !c.IsActive() {
		t.Error("expected reconstituted requester to default to active")
	}
}

func TestRequester_SetStatus(t *testing.T) {
	cases := []struct {
		name     string
		status   RequesterStatus
		expected RequesterStatus
	}{
		{"inativa", RequesterStatusInactive, RequesterStatusInactive},
		{"reativa", RequesterStatusActive, RequesterStatusActive},
		{"status desconhecido e ignorado", RequesterStatus("banido"), RequesterStatusActive},
		{"status vazio e ignorado", RequesterStatus(""), RequesterStatusActive},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewRequester("João", "52998224725", "j@j.com", "11999")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			c.SetStatus(tc.status)

			if c.Status() != tc.expected {
				t.Errorf("expected status %q, got %q", tc.expected, c.Status())
			}
		})
	}
}

func TestRequester_SetStatus_DoesNotTouchUpdatedAt(t *testing.T) {
	c, err := NewRequester("João", "52998224725", "j@j.com", "11999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	before := c.UpdatedAt()
	time.Sleep(time.Millisecond)
	c.SetStatus(RequesterStatusInactive)

	if !c.UpdatedAt().Equal(before) {
		t.Error("expected SetStatus to preserve UpdatedAt during reconstitution")
	}
}

func TestRequesterStatus_IsValid(t *testing.T) {
	valid := []RequesterStatus{RequesterStatusActive, RequesterStatusInactive}
	for _, s := range valid {
		if !s.IsValid() {
			t.Errorf("expected %q to be valid", s)
		}
	}

	if RequesterStatus("qualquer").IsValid() {
		t.Error("expected unknown status to be invalid")
	}
}
