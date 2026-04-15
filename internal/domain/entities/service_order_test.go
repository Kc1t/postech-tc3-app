package entities

import (
	"testing"
	"time"

	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
)

func TestUpdateStatus_TransicoesValidas(t *testing.T) {
	tests := []struct {
		name string
		from OrderStatus
		to   OrderStatus
	}{
		{"received -> in_diagnosis", StatusReceived, StatusInDiagnosis},
		{"in_diagnosis -> awaiting_approval", StatusInDiagnosis, StatusAwaitingApproval},
		{"awaiting_approval -> in_execution", StatusAwaitingApproval, StatusInExecution},
		{"awaiting_approval -> received (recusa)", StatusAwaitingApproval, StatusReceived},
		{"in_execution -> finished", StatusInExecution, StatusFinished},
		{"finished -> delivered", StatusFinished, StatusDelivered},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			so := newServiceOrderWithStatus(tt.from)
			err := so.UpdateStatus(tt.to)
			if err != nil {
				t.Fatalf("esperava sucesso, obteve erro: %v", err)
			}
			if so.Status() != tt.to {
				t.Errorf("Status() = %q, esperava %q", so.Status(), tt.to)
			}
		})
	}
}

func TestUpdateStatus_TransicoesInvalidas(t *testing.T) {
	tests := []struct {
		name string
		from OrderStatus
		to   OrderStatus
	}{
		{"received -> finished (pular etapas)", StatusReceived, StatusFinished},
		{"received -> delivered (pular tudo)", StatusReceived, StatusDelivered},
		{"received -> in_execution (pular diag)", StatusReceived, StatusInExecution},
		{"in_diagnosis -> in_execution (pular aprov)", StatusInDiagnosis, StatusInExecution},
		{"in_execution -> delivered (pular finished)", StatusInExecution, StatusDelivered},
		{"delivered -> received (estado terminal)", StatusDelivered, StatusReceived},
		{"delivered -> in_diagnosis (estado terminal)", StatusDelivered, StatusInDiagnosis},
		{"finished -> in_execution (voltar atras)", StatusFinished, StatusInExecution},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			so := newServiceOrderWithStatus(tt.from)
			err := so.UpdateStatus(tt.to)
			if err == nil {
				t.Fatal("esperava erro, obteve sucesso")
			}
			if err != domainerrors.ErrInvalidStatus {
				t.Errorf("erro = %v, esperava ErrInvalidStatus", err)
			}
			// Status nao deve ter mudado
			if so.Status() != tt.from {
				t.Errorf("Status() = %q, deveria permanecer %q", so.Status(), tt.from)
			}
		})
	}
}

func TestUpdateStatus_StatusDesconhecido(t *testing.T) {
	so := NewServiceOrder("cust-1", "veh-1")
	err := so.UpdateStatus(OrderStatus("inexistente"))
	if err == nil {
		t.Fatal("esperava erro para status desconhecido")
	}
	if err != domainerrors.ErrInvalidStatusValue {
		t.Errorf("erro = %v, esperava ErrInvalidStatusValue", err)
	}
}

func TestUpdateStatus_NaoAlteraQuandoInvalido(t *testing.T) {
	so := NewServiceOrder("cust-1", "veh-1")
	original := so.Status()
	_ = so.UpdateStatus(StatusFinished) // transicao invalida
	if so.Status() != original {
		t.Errorf("Status mudou para %q apos transicao invalida, deveria ser %q", so.Status(), original)
	}
}

func TestNewServiceOrder_StatusInicial(t *testing.T) {
	so := NewServiceOrder("cust-1", "veh-1")
	if so.Status() != StatusReceived {
		t.Errorf("Status() = %q, esperava %q", so.Status(), StatusReceived)
	}
}

func TestIsValidStatus(t *testing.T) {
	validos := []OrderStatus{
		StatusReceived, StatusInDiagnosis, StatusAwaitingApproval,
		StatusInExecution, StatusFinished, StatusDelivered,
	}
	for _, s := range validos {
		if !IsValidStatus(s) {
			t.Errorf("IsValidStatus(%q) = false, esperava true", s)
		}
	}

	invalidos := []OrderStatus{"", "cancelado", "pendente", "xyz"}
	for _, s := range invalidos {
		if IsValidStatus(s) {
			t.Errorf("IsValidStatus(%q) = true, esperava false", s)
		}
	}
}

func TestFluxoCompleto_ReceivedAteDelivered(t *testing.T) {
	so := NewServiceOrder("cust-1", "veh-1")
	fluxo := []OrderStatus{
		StatusInDiagnosis,
		StatusAwaitingApproval,
		StatusInExecution,
		StatusFinished,
		StatusDelivered,
	}

	for _, status := range fluxo {
		if err := so.UpdateStatus(status); err != nil {
			t.Fatalf("falha na transicao para %q: %v", status, err)
		}
	}

	if so.Status() != StatusDelivered {
		t.Errorf("Status final = %q, esperava %q", so.Status(), StatusDelivered)
	}
}

func TestFluxoRecusa_VoltaParaReceived(t *testing.T) {
	so := NewServiceOrder("cust-1", "veh-1")

	// received -> in_diagnosis -> awaiting_approval -> received (recusa)
	if err := so.UpdateStatus(StatusInDiagnosis); err != nil {
		t.Fatalf("transicao para in_diagnosis falhou: %v", err)
	}
	if err := so.UpdateStatus(StatusAwaitingApproval); err != nil {
		t.Fatalf("transicao para awaiting_approval falhou: %v", err)
	}

	if err := so.UpdateStatus(StatusReceived); err != nil {
		t.Fatalf("recusa deveria ser permitida: %v", err)
	}
	if so.Status() != StatusReceived {
		t.Errorf("Status() = %q, esperava %q", so.Status(), StatusReceived)
	}

	// Deve poder reiniciar o fluxo normalmente
	if err := so.UpdateStatus(StatusInDiagnosis); err != nil {
		t.Fatalf("reinicio do fluxo deveria ser permitido: %v", err)
	}
}

func TestRecalcTotal_ServicosEPecas(t *testing.T) {
	so := NewServiceOrder("cust-1", "veh-1")

	so.AddService(ServiceItem{ServiceID: "s1", Description: "Troca de oleo", Price: 150.0})
	so.AddService(ServiceItem{ServiceID: "s2", Description: "Alinhamento", Price: 80.0})
	so.AddPart(PartItem{PartID: "p1", Description: "Filtro", Quantity: 2, UnitPrice: 25.0})

	// 150 + 80 + (2 * 25) = 280
	expected := 280.0
	if so.TotalAmount() != expected {
		t.Errorf("TotalAmount() = %.2f, esperava %.2f", so.TotalAmount(), expected)
	}
}

func TestSetServices_SubstituiERecalcula(t *testing.T) {
	so := NewServiceOrder("cust-1", "veh-1")
	so.AddService(ServiceItem{ServiceID: "s1", Price: 100.0})

	so.SetServices([]ServiceItem{
		{ServiceID: "s2", Price: 200.0},
		{ServiceID: "s3", Price: 300.0},
	})

	if so.TotalAmount() != 500.0 {
		t.Errorf("TotalAmount() = %.2f, esperava 500.00", so.TotalAmount())
	}
	if len(so.Services()) != 2 {
		t.Errorf("len(Services()) = %d, esperava 2", len(so.Services()))
	}
}

func TestSetParts_SubstituiERecalcula(t *testing.T) {
	so := NewServiceOrder("cust-1", "veh-1")
	so.AddService(ServiceItem{ServiceID: "s1", Price: 100.0})
	so.AddPart(PartItem{PartID: "p1", Quantity: 1, UnitPrice: 50.0})

	so.SetParts([]PartItem{
		{PartID: "p2", Quantity: 3, UnitPrice: 20.0},
	})

	// 100 (servico) + 3*20 (peca) = 160
	if so.TotalAmount() != 160.0 {
		t.Errorf("TotalAmount() = %.2f, esperava 160.00", so.TotalAmount())
	}
}

// Helper para criar OS com status especifico (via reconstituicao).
func newServiceOrderWithStatus(status OrderStatus) *ServiceOrder {
	t := time.Now()
	return ReconstituteServiceOrder(
		"order-1", "cust-1", "veh-1",
		status,
		nil, nil, 0, "",
		t, t,
	)
}
