package entities

import (
	"testing"
	"time"
)

func TestNewServiceOrder(t *testing.T) {
	before := time.Now()
	so := NewServiceOrder("cust-1", "veh-1")
	after := time.Now()

	if so.CustomerID() != "cust-1" {
		t.Errorf("expected customerID %q, got %q", "cust-1", so.CustomerID())
	}
	if so.VehicleID() != "veh-1" {
		t.Errorf("expected vehicleID %q, got %q", "veh-1", so.VehicleID())
	}
	if so.Status() != StatusReceived {
		t.Errorf("expected status %q, got %q", StatusReceived, so.Status())
	}
	if so.TotalAmount() != 0 {
		t.Errorf("expected totalAmount 0, got %v", so.TotalAmount())
	}
	if len(so.Services()) != 0 {
		t.Errorf("expected empty services, got %d", len(so.Services()))
	}
	if len(so.Parts()) != 0 {
		t.Errorf("expected empty parts, got %d", len(so.Parts()))
	}
	if so.ID() != "" {
		t.Errorf("expected empty ID, got %q", so.ID())
	}
	if so.CreatedAt().Before(before) || so.CreatedAt().After(after) {
		t.Error("CreatedAt out of expected range")
	}
}

func TestReconstituteServiceOrder(t *testing.T) {
	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	services := []ServiceItem{{ServiceID: "s1", Price: 100.0}}
	parts := []PartItem{{PartID: "p1", Quantity: 2, UnitPrice: 50.0}}

	so := ReconstituteServiceOrder("oid-1", "cust-1", "veh-1", StatusInExecution, services, parts, 200.0, "notas", createdAt, updatedAt)

	if so.ID() != "oid-1" {
		t.Errorf("expected ID %q, got %q", "oid-1", so.ID())
	}
	if so.Status() != StatusInExecution {
		t.Errorf("expected status %q, got %q", StatusInExecution, so.Status())
	}
	if so.TotalAmount() != 200.0 {
		t.Errorf("expected totalAmount 200.0, got %v", so.TotalAmount())
	}
	if so.Notes() != "notas" {
		t.Errorf("expected notes %q, got %q", "notas", so.Notes())
	}
	if len(so.Services()) != 1 {
		t.Errorf("expected 1 service, got %d", len(so.Services()))
	}
	if len(so.Parts()) != 1 {
		t.Errorf("expected 1 part, got %d", len(so.Parts()))
	}
}

func TestServiceOrder_SetID(t *testing.T) {
	so := NewServiceOrder("cust-1", "veh-1")
	so.SetID("oid-1")
	if so.ID() != "oid-1" {
		t.Errorf("expected ID %q, got %q", "oid-1", so.ID())
	}
}

func TestServiceOrder_SetNotes_TouchesUpdatedAt(t *testing.T) {
	so := NewServiceOrder("cust-1", "veh-1")
	before := so.UpdatedAt()
	time.Sleep(time.Millisecond)
	so.SetNotes("Trocar pastilhas também")

	if so.Notes() != "Trocar pastilhas também" {
		t.Errorf("expected notes %q, got %q", "Trocar pastilhas também", so.Notes())
	}
	if !so.UpdatedAt().After(before) {
		t.Error("expected UpdatedAt to be updated after SetNotes")
	}
}

func TestServiceOrder_UpdateStatus_TouchesUpdatedAt(t *testing.T) {
	so := NewServiceOrder("cust-1", "veh-1")
	before := so.UpdatedAt()
	time.Sleep(time.Millisecond)
	so.UpdateStatus(StatusInDiagnosis)

	if so.Status() != StatusInDiagnosis {
		t.Errorf("expected status %q, got %q", StatusInDiagnosis, so.Status())
	}
	if !so.UpdatedAt().After(before) {
		t.Error("expected UpdatedAt to be updated after UpdateStatus")
	}
}

func TestServiceOrder_AddService_RecalcTotal(t *testing.T) {
	so := NewServiceOrder("cust-1", "veh-1")
	so.AddService(ServiceItem{ServiceID: "s1", Description: "Troca de óleo", Price: 150.0})

	if len(so.Services()) != 1 {
		t.Fatalf("expected 1 service, got %d", len(so.Services()))
	}
	if so.TotalAmount() != 150.0 {
		t.Errorf("expected totalAmount 150.0, got %v", so.TotalAmount())
	}
}

func TestServiceOrder_AddPart_RecalcTotal(t *testing.T) {
	so := NewServiceOrder("cust-1", "veh-1")
	so.AddPart(PartItem{PartID: "p1", Description: "Filtro", Quantity: 2, UnitPrice: 50.0})

	if len(so.Parts()) != 1 {
		t.Fatalf("expected 1 part, got %d", len(so.Parts()))
	}
	if so.TotalAmount() != 100.0 {
		t.Errorf("expected totalAmount 100.0, got %v", so.TotalAmount())
	}
}

func TestServiceOrder_AddServiceAndPart_RecalcTotal(t *testing.T) {
	so := NewServiceOrder("cust-1", "veh-1")
	so.AddService(ServiceItem{ServiceID: "s1", Price: 200.0})
	so.AddPart(PartItem{PartID: "p1", Quantity: 3, UnitPrice: 30.0})

	expected := 200.0 + 3*30.0
	if so.TotalAmount() != expected {
		t.Errorf("expected totalAmount %v, got %v", expected, so.TotalAmount())
	}
}

func TestServiceOrder_MultipleServices_RecalcTotal(t *testing.T) {
	so := NewServiceOrder("cust-1", "veh-1")
	so.AddService(ServiceItem{ServiceID: "s1", Price: 100.0})
	so.AddService(ServiceItem{ServiceID: "s2", Price: 200.0})

	if len(so.Services()) != 2 {
		t.Fatalf("expected 2 services, got %d", len(so.Services()))
	}
	if so.TotalAmount() != 300.0 {
		t.Errorf("expected totalAmount 300.0, got %v", so.TotalAmount())
	}
}
