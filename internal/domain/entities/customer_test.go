package entities

import (
	"testing"
	"time"
)

func TestNewCustomer(t *testing.T) {
	before := time.Now()
	c := NewCustomer("João Silva", "12345678901", "joao@email.com", "11999990000")
	after := time.Now()

	if c.Name() != "João Silva" {
		t.Errorf("expected name %q, got %q", "João Silva", c.Name())
	}
	if c.Document() != "12345678901" {
		t.Errorf("expected document %q, got %q", "12345678901", c.Document())
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

func TestReconstituteCustomer(t *testing.T) {
	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	c := ReconstituteCustomer("id-1", "Maria", "98765432100", "maria@email.com", "11888880000", createdAt, updatedAt)

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

func TestCustomer_SetID(t *testing.T) {
	c := NewCustomer("João", "123", "j@j.com", "11999")
	c.SetID("new-id")
	if c.ID() != "new-id" {
		t.Errorf("expected ID %q, got %q", "new-id", c.ID())
	}
}

func TestCustomer_SetName_TouchesUpdatedAt(t *testing.T) {
	c := NewCustomer("João", "123", "j@j.com", "11999")
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

func TestCustomer_SetEmail_TouchesUpdatedAt(t *testing.T) {
	c := NewCustomer("João", "123", "j@j.com", "11999")
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

func TestCustomer_SetPhone_TouchesUpdatedAt(t *testing.T) {
	c := NewCustomer("João", "123", "j@j.com", "11999")
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
