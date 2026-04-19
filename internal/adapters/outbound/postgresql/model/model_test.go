package pgmodel

import (
	"testing"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
)

// --- Customer ---

func TestFromCustomer_ToDomain_RoundTrip(t *testing.T) {
	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	c := entities.ReconstituteCustomer("id-1", "João", "12345678901", "j@j.com", "11999", createdAt, updatedAt)

	m := FromCustomer(c)
	if m.ID != "id-1" || m.Name != "João" || m.Document != "12345678901" {
		t.Fatalf("FromCustomer fields mismatch: %+v", m)
	}

	back := m.ToDomain()
	if back.ID() != "id-1" || back.Name() != "João" || back.Document() != "12345678901" {
		t.Fatalf("ToDomain fields mismatch: %+v", back)
	}
	if !back.CreatedAt().Equal(createdAt) {
		t.Errorf("CreatedAt mismatch: %v != %v", back.CreatedAt(), createdAt)
	}
}

// --- Vehicle ---

func TestFromVehicle_ToDomain_RoundTrip(t *testing.T) {
	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	v := entities.ReconstituteVehicle("vid-1", "cust-1", "ABC1234", "Toyota", "Corolla", 2020, createdAt, updatedAt)

	m := FromVehicle(v)
	if m.ID != "vid-1" || m.Plate != "ABC1234" || m.Year != 2020 {
		t.Fatalf("FromVehicle fields mismatch: %+v", m)
	}

	back := m.ToDomain()
	if back.ID() != "vid-1" || back.Plate() != "ABC1234" || back.Year() != 2020 {
		t.Fatalf("ToDomain fields mismatch: %+v", back)
	}
}

// --- Part ---

func TestFromPart_ToDomain_RoundTrip(t *testing.T) {
	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	p := entities.ReconstitutePart("pid-1", "Filtro", "Desc", "unidade", 49.90, 10, createdAt, updatedAt)

	m := FromPart(p)
	if m.ID != "pid-1" || m.Name != "Filtro" || m.Stock != 10 {
		t.Fatalf("FromPart fields mismatch: %+v", m)
	}

	back := m.ToDomain()
	if back.ID() != "pid-1" || back.Name() != "Filtro" || back.Stock() != 10 {
		t.Fatalf("ToDomain fields mismatch: %+v", back)
	}
}

// --- Service ---

func TestFromService_ToDomain_RoundTrip(t *testing.T) {
	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	s := entities.ReconstituteService("sid-1", "Troca de óleo", "Desc", 150.0, 60, createdAt, updatedAt)

	m := FromService(s)
	if m.ID != "sid-1" || m.Name != "Troca de óleo" || m.DurationMin != 60 {
		t.Fatalf("FromService fields mismatch: %+v", m)
	}

	back := m.ToDomain()
	if back.ID() != "sid-1" || back.Name() != "Troca de óleo" || back.DurationMin() != 60 {
		t.Fatalf("ToDomain fields mismatch: %+v", back)
	}
}

// --- ServiceOrder ---

func TestFromServiceOrder_ToDomain_RoundTrip(t *testing.T) {
	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	services := []entities.ServiceItem{{ServiceID: "s1", Description: "Desc", Price: 100.0}}
	parts := []entities.PartItem{{PartID: "p1", Description: "Peca", Quantity: 2, UnitPrice: 50.0}}
	so := entities.ReconstituteServiceOrder(
		"oid-1", "cust-1", "veh-1", entities.StatusInExecution,
		services, parts, 200.0, "notas", createdAt, updatedAt,
	)

	m := FromServiceOrder(so)
	if m.ID != "oid-1" || m.Status != entities.StatusInExecution {
		t.Fatalf("FromServiceOrder fields mismatch: %+v", m)
	}
	if len(m.Services) != 1 || len(m.Parts) != 1 {
		t.Fatalf("expected 1 service and 1 part, got %d and %d", len(m.Services), len(m.Parts))
	}

	back := m.ToDomain()
	if back.ID() != "oid-1" || back.Status() != entities.StatusInExecution {
		t.Fatalf("ToDomain fields mismatch: %+v", back)
	}
	if len(back.Services()) != 1 || len(back.Parts()) != 1 {
		t.Fatalf("expected 1 service and 1 part back, got %d and %d", len(back.Services()), len(back.Parts()))
	}
	if back.TotalAmount() != 200.0 {
		t.Errorf("expected totalAmount 200.0, got %v", back.TotalAmount())
	}
}

func TestFromServiceOrder_EmptyItemsSlices(t *testing.T) {
	so := entities.NewServiceOrder("cust-1", "veh-1")
	m := FromServiceOrder(so)
	if len(m.Services) != 0 || len(m.Parts) != 0 {
		t.Fatalf("expected empty slices, got %d services and %d parts", len(m.Services), len(m.Parts))
	}
	back := m.ToDomain()
	if len(back.Services()) != 0 || len(back.Parts()) != 0 {
		t.Fatalf("expected empty back slices")
	}
}
