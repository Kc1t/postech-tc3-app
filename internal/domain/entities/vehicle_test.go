package entities

import (
	"testing"
	"time"
)

func TestNewVehicle(t *testing.T) {
	before := time.Now()
	v, err := NewVehicle("cust-1", "ABC1234", "Toyota", "Corolla", 2020)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	after := time.Now()

	if v.RequesterID() != "cust-1" {
		t.Errorf("expected requesterID %q, got %q", "cust-1", v.RequesterID())
	}
	if v.Plate() != "ABC1234" {
		t.Errorf("expected plate %q, got %q", "ABC1234", v.Plate())
	}
	if v.Brand() != "Toyota" {
		t.Errorf("expected brand %q, got %q", "Toyota", v.Brand())
	}
	if v.Model() != "Corolla" {
		t.Errorf("expected model %q, got %q", "Corolla", v.Model())
	}
	if v.Year() != 2020 {
		t.Errorf("expected year %d, got %d", 2020, v.Year())
	}
	if v.ID() != "" {
		t.Errorf("expected empty ID, got %q", v.ID())
	}
	if v.CreatedAt().Before(before) || v.CreatedAt().After(after) {
		t.Error("CreatedAt out of expected range")
	}
}

func TestNewVehicle_InvalidPlate(t *testing.T) {
	_, err := NewVehicle("cust-1", "INVALID", "Toyota", "Corolla", 2020)
	if err == nil {
		t.Fatal("expected error for invalid plate, got nil")
	}
}

func TestReconstituteVehicle(t *testing.T) {
	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	v := ReconstituteVehicle("id-v", "cust-1", "XYZ9876", "Honda", "Civic", 2019, createdAt, updatedAt)

	if v.ID() != "id-v" {
		t.Errorf("expected ID %q, got %q", "id-v", v.ID())
	}
	if v.Plate() != "XYZ9876" {
		t.Errorf("expected plate %q, got %q", "XYZ9876", v.Plate())
	}
	if !v.CreatedAt().Equal(createdAt) {
		t.Errorf("expected createdAt %v, got %v", createdAt, v.CreatedAt())
	}
}

func TestVehicle_SetID(t *testing.T) {
	v, err := NewVehicle("cust-1", "ABC1234", "Toyota", "Corolla", 2020)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v.SetID("vid-1")
	if v.ID() != "vid-1" {
		t.Errorf("expected ID %q, got %q", "vid-1", v.ID())
	}
}

func TestVehicle_SetRequesterID(t *testing.T) {
	v, err := NewVehicle("cust-1", "ABC1234", "Toyota", "Corolla", 2020)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v.SetRequesterID("cust-2")
	if v.RequesterID() != "cust-2" {
		t.Errorf("expected requesterID %q, got %q", "cust-2", v.RequesterID())
	}
}

func TestVehicle_SetPlate_TouchesUpdatedAt(t *testing.T) {
	v, err := NewVehicle("cust-1", "ABC1234", "Toyota", "Corolla", 2020)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	before := v.UpdatedAt()
	time.Sleep(time.Millisecond)
	if err := v.SetPlate("XYZ5678"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v.Plate() != "XYZ5678" {
		t.Errorf("expected plate %q, got %q", "XYZ5678", v.Plate())
	}
	if !v.UpdatedAt().After(before) {
		t.Error("expected UpdatedAt to be updated after SetPlate")
	}
}

func TestVehicle_SetBrand_TouchesUpdatedAt(t *testing.T) {
	v, err := NewVehicle("cust-1", "ABC1234", "Toyota", "Corolla", 2020)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	before := v.UpdatedAt()
	time.Sleep(time.Millisecond)
	v.SetBrand("Honda")

	if v.Brand() != "Honda" {
		t.Errorf("expected brand %q, got %q", "Honda", v.Brand())
	}
	if !v.UpdatedAt().After(before) {
		t.Error("expected UpdatedAt to be updated after SetBrand")
	}
}

func TestVehicle_SetModel_TouchesUpdatedAt(t *testing.T) {
	v, err := NewVehicle("cust-1", "ABC1234", "Toyota", "Corolla", 2020)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	before := v.UpdatedAt()
	time.Sleep(time.Millisecond)
	v.SetModel("Fit")

	if v.Model() != "Fit" {
		t.Errorf("expected model %q, got %q", "Fit", v.Model())
	}
	if !v.UpdatedAt().After(before) {
		t.Error("expected UpdatedAt to be updated after SetModel")
	}
}

func TestVehicle_SetYear_TouchesUpdatedAt(t *testing.T) {
	v, err := NewVehicle("cust-1", "ABC1234", "Toyota", "Corolla", 2020)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	before := v.UpdatedAt()
	time.Sleep(time.Millisecond)
	v.SetYear(2022)

	if v.Year() != 2022 {
		t.Errorf("expected year %d, got %d", 2022, v.Year())
	}
	if !v.UpdatedAt().After(before) {
		t.Error("expected UpdatedAt to be updated after SetYear")
	}
}
