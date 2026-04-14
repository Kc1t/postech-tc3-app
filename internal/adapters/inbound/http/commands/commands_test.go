package commands

import (
	"testing"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
)

// --- Customer ---

func TestCreateCustomerRequest_ToDomain(t *testing.T) {
	req := CreateCustomerRequest{Name: "João", Document: "52998224725", Email: "j@j.com", Phone: "11999"}
	c, err := req.ToDomain()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Name() != "João" || c.Document() != "52998224725" || c.Email() != "j@j.com" || c.Phone() != "11999" {
		t.Fatalf("unexpected domain customer: %+v", c)
	}
}

func TestToCustomerResponse(t *testing.T) {
	now := time.Now()
	c := entities.ReconstituteCustomer("id-1", "João", "123", "j@j.com", "11999", now, now)
	resp := ToCustomerResponse(c)
	if resp.Name != "João" || resp.Document != "123" || resp.Email != "j@j.com" {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if _, err := time.Parse(time.RFC3339, resp.CreatedAt); err != nil {
		t.Errorf("invalid CreatedAt format: %q", resp.CreatedAt)
	}
}

func TestToCustomerListResponse(t *testing.T) {
	now := time.Now()
	customers := []*entities.Customer{
		entities.ReconstituteCustomer("", "João", "123", "j@j.com", "11999", now, now),
		entities.ReconstituteCustomer("", "Maria", "456", "m@m.com", "11888", now, now),
	}
	resp := ToCustomerListResponse(customers)
	if len(resp) != 2 {
		t.Fatalf("expected 2, got %d", len(resp))
	}
}

// --- Vehicle ---

func TestCreateVehicleRequest_ToDomain(t *testing.T) {
	req := CreateVehicleRequest{CustomerID: "cust-1", Plate: "ABC1234", Brand: "Toyota", Model: "Corolla", Year: 2020}
	v, err := req.ToDomain()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Plate() != "ABC1234" || v.Brand() != "Toyota" || v.Year() != 2020 {
		t.Fatalf("unexpected domain vehicle: %+v", v)
	}
}

func TestToVehicleResponse(t *testing.T) {
	now := time.Now()
	v := entities.ReconstituteVehicle("", "cust-1", "ABC1234", "Toyota", "Corolla", 2020, now, now)
	resp := ToVehicleResponse(v)
	if resp.Plate != "ABC1234" || resp.Year != 2020 {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if _, err := time.Parse(time.RFC3339, resp.CreatedAt); err != nil {
		t.Errorf("invalid CreatedAt format: %q", resp.CreatedAt)
	}
}

func TestToVehicleListResponse(t *testing.T) {
	now := time.Now()
	vehicles := []*entities.Vehicle{
		entities.ReconstituteVehicle("", "cust-1", "ABC1234", "Toyota", "Corolla", 2020, now, now),
	}
	resp := ToVehicleListResponse(vehicles)
	if len(resp) != 1 {
		t.Fatalf("expected 1, got %d", len(resp))
	}
}

// --- Part ---

func TestCreatePartRequest_ToDomain(t *testing.T) {
	req := CreatePartRequest{Name: "Filtro", Unit: "unidade", Price: 49.90, Stock: 10}
	p := req.ToDomain()
	if p.Name() != "Filtro" || p.Unit() != "unidade" || p.Price() != 49.90 || p.Stock() != 10 {
		t.Fatalf("unexpected domain part: %+v", p)
	}
}

func TestToPartResponse(t *testing.T) {
	p := entities.NewPart("Filtro", "Desc", "unidade", 49.90, 10)
	resp := ToPartResponse(p)
	if resp.Name != "Filtro" || resp.Price != 49.90 || resp.Stock != 10 {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if _, err := time.Parse(time.RFC3339, resp.CreatedAt); err != nil {
		t.Errorf("invalid CreatedAt format: %q", resp.CreatedAt)
	}
}

func TestToPartListResponse(t *testing.T) {
	parts := []*entities.Part{entities.NewPart("Filtro", "Desc", "unidade", 49.90, 10)}
	resp := ToPartListResponse(parts)
	if len(resp) != 1 {
		t.Fatalf("expected 1, got %d", len(resp))
	}
}

// --- Service ---

func TestCreateServiceRequest_ToDomain(t *testing.T) {
	req := CreateServiceRequest{Name: "Troca de óleo", Price: 150.0, DurationMin: 60}
	s := req.ToDomain()
	if s.Name() != "Troca de óleo" || s.Price() != 150.0 || s.DurationMin() != 60 {
		t.Fatalf("unexpected domain service: %+v", s)
	}
}

func TestToServiceResponse(t *testing.T) {
	s := entities.NewService("Troca de óleo", "Desc", 150.0, 60)
	resp := ToServiceResponse(s)
	if resp.Name != "Troca de óleo" || resp.Price != 150.0 || resp.DurationMin != 60 {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if _, err := time.Parse(time.RFC3339, resp.CreatedAt); err != nil {
		t.Errorf("invalid CreatedAt format: %q", resp.CreatedAt)
	}
}

func TestToServiceListResponse(t *testing.T) {
	svcs := []*entities.Service{entities.NewService("Troca", "Desc", 100.0, 30)}
	resp := ToServiceListResponse(svcs)
	if len(resp) != 1 {
		t.Fatalf("expected 1, got %d", len(resp))
	}
}

// --- ServiceOrder ---

func TestCreateServiceOrderRequest_ToDomain_WithItems(t *testing.T) {
	req := CreateServiceOrderRequest{
		CustomerID: "cust-1",
		VehicleID:  "veh-1",
		Notes:      "trocar pastilhas",
		Services: []ServiceItemRequest{
			{ServiceID: "s1", Description: "Troca", Price: 100.0},
		},
		Parts: []PartItemRequest{
			{PartID: "p1", Description: "Filtro", Quantity: 2, UnitPrice: 50.0},
		},
	}
	so := req.ToDomain()
	if so.CustomerID() != "cust-1" || so.VehicleID() != "veh-1" {
		t.Fatalf("unexpected customerID/vehicleID: %s/%s", so.CustomerID(), so.VehicleID())
	}
	if so.Notes() != "trocar pastilhas" {
		t.Errorf("expected notes %q, got %q", "trocar pastilhas", so.Notes())
	}
	if len(so.Services()) != 1 || len(so.Parts()) != 1 {
		t.Fatalf("expected 1 service and 1 part, got %d and %d", len(so.Services()), len(so.Parts()))
	}
}

func TestCreateServiceOrderRequest_ToDomain_Empty(t *testing.T) {
	req := CreateServiceOrderRequest{CustomerID: "cust-1", VehicleID: "veh-1"}
	so := req.ToDomain()
	if len(so.Services()) != 0 || len(so.Parts()) != 0 {
		t.Fatalf("expected empty services/parts")
	}
}

func TestToServiceOrderResponse(t *testing.T) {
	so := entities.NewServiceOrder("cust-1", "veh-1")
	so.AddService(entities.ServiceItem{ServiceID: "s1", Description: "Desc", Price: 100.0})
	so.AddPart(entities.PartItem{PartID: "p1", Description: "Peca", Quantity: 2, UnitPrice: 50.0})

	resp := ToServiceOrderResponse(so)
	if resp.CustomerID != "cust-1" || resp.VehicleID != "veh-1" {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if len(resp.Services) != 1 || len(resp.Parts) != 1 {
		t.Fatalf("expected 1 service and 1 part, got %d and %d", len(resp.Services), len(resp.Parts))
	}
	if resp.TotalAmount != 200.0 {
		t.Errorf("expected totalAmount 200.0, got %v", resp.TotalAmount)
	}
}

func TestToServiceOrderListResponse(t *testing.T) {
	orders := []*entities.ServiceOrder{
		entities.NewServiceOrder("cust-1", "veh-1"),
		entities.NewServiceOrder("cust-2", "veh-2"),
	}
	resp := ToServiceOrderListResponse(orders)
	if len(resp) != 2 {
		t.Fatalf("expected 2, got %d", len(resp))
	}
}
