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

func TestAuthorizeRequesterTransition_StatusPermitidos(t *testing.T) {
	tests := []struct {
		name string
		to   OrderStatus
	}{
		{"aprovar (awaiting_approval -> in_execution)", StatusInExecution},
		{"rejeitar (awaiting_approval -> received)", StatusReceived},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			so := newServiceOrderWithStatus(StatusAwaitingApproval)
			if err := so.AuthorizeRequesterTransition(tt.to); err != nil {
				t.Fatalf("esperava sucesso, obteve erro: %v", err)
			}
			if so.Status() != tt.to {
				t.Errorf("Status() = %q, esperava %q", so.Status(), tt.to)
			}
		})
	}
}

func TestAuthorizeRequesterTransition_StatusProibidos(t *testing.T) {
	// Status que so o mecanico/atendente (rota autenticada) pode disparar.
	tests := []struct {
		name string
		from OrderStatus
		to   OrderStatus
	}{
		{"received -> in_diagnosis", StatusReceived, StatusInDiagnosis},
		{"in_diagnosis -> awaiting_approval", StatusInDiagnosis, StatusAwaitingApproval},
		{"in_execution -> finished", StatusInExecution, StatusFinished},
		{"finished -> delivered", StatusFinished, StatusDelivered},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			so := newServiceOrderWithStatus(tt.from)
			err := so.AuthorizeRequesterTransition(tt.to)
			if err != domainerrors.ErrStatusNotAllowedForRequester {
				t.Fatalf("erro = %v, esperava ErrStatusNotAllowedForRequester", err)
			}
			if so.Status() != tt.from {
				t.Errorf("Status() = %q, deveria permanecer %q", so.Status(), tt.from)
			}
		})
	}
}

func TestAuthorizeRequesterTransition_WhitelistPassaMasMaquinaReprova(t *testing.T) {
	// in_execution esta na whitelist, mas a transicao received -> in_execution
	// nao existe na maquina de estados: deve retornar ErrInvalidStatus.
	so := newServiceOrderWithStatus(StatusReceived)
	err := so.AuthorizeRequesterTransition(StatusInExecution)
	if err != domainerrors.ErrInvalidStatus {
		t.Fatalf("erro = %v, esperava ErrInvalidStatus", err)
	}
	if so.Status() != StatusReceived {
		t.Errorf("Status() = %q, deveria permanecer %q", so.Status(), StatusReceived)
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

// =============================================================================
// Testes dos timestamps de ciclo de vida (startedAt, finishedAt)
// =============================================================================

func TestUpdateStatus_InExecution_GravaStartedAt(t *testing.T) {
	so := newServiceOrderWithStatus(StatusAwaitingApproval)

	if so.StartedAt() != nil {
		t.Fatal("StartedAt() deveria ser nil antes da aprovacao")
	}
	before := time.Now()
	if err := so.UpdateStatus(StatusInExecution); err != nil {
		t.Fatalf("transicao falhou: %v", err)
	}
	after := time.Now()

	if so.StartedAt() == nil {
		t.Fatal("StartedAt() nao foi gravado ao entrar em in_execution")
	}
	if so.StartedAt().Before(before) || so.StartedAt().After(after) {
		t.Errorf("StartedAt() = %v, esperava entre %v e %v", *so.StartedAt(), before, after)
	}
	if so.FinishedAt() != nil {
		t.Error("FinishedAt() deveria continuar nil")
	}
}

func TestUpdateStatus_Finished_GravaFinishedAt(t *testing.T) {
	so := newServiceOrderWithStatus(StatusInExecution)

	if so.FinishedAt() != nil {
		t.Fatal("FinishedAt() deveria ser nil antes de finalizar")
	}
	before := time.Now()
	if err := so.UpdateStatus(StatusFinished); err != nil {
		t.Fatalf("transicao falhou: %v", err)
	}
	after := time.Now()

	if so.FinishedAt() == nil {
		t.Fatal("FinishedAt() nao foi gravado ao entrar em finished")
	}
	if so.FinishedAt().Before(before) || so.FinishedAt().After(after) {
		t.Errorf("FinishedAt() = %v, esperava entre %v e %v", *so.FinishedAt(), before, after)
	}
}

// Transicoes que nao passam por in_execution/finished nao devem tocar os timestamps.
func TestUpdateStatus_OutrasTransicoes_NaoTocamTimestamps(t *testing.T) {
	tests := []struct {
		from OrderStatus
		to   OrderStatus
	}{
		{StatusReceived, StatusInDiagnosis},
		{StatusInDiagnosis, StatusAwaitingApproval},
		{StatusAwaitingApproval, StatusReceived}, // recusa
		{StatusFinished, StatusDelivered},
	}
	for _, tt := range tests {
		t.Run(string(tt.from)+" -> "+string(tt.to), func(t *testing.T) {
			so := newServiceOrderWithStatus(tt.from)
			if err := so.UpdateStatus(tt.to); err != nil {
				t.Fatalf("transicao falhou: %v", err)
			}
			if so.StartedAt() != nil {
				t.Errorf("StartedAt() = %v, deveria continuar nil", so.StartedAt())
			}
			if so.FinishedAt() != nil {
				t.Errorf("FinishedAt() = %v, deveria continuar nil", so.FinishedAt())
			}
		})
	}
}

// Fluxo completo: startedAt e finishedAt sao gravados uma unica vez, no
// momento correto do ciclo de vida.
func TestFluxoCompleto_PopulaStartedEFinished(t *testing.T) {
	so := NewServiceOrder("cust-1", "veh-1")
	_ = so.UpdateStatus(StatusInDiagnosis)
	_ = so.UpdateStatus(StatusAwaitingApproval)

	if so.StartedAt() != nil || so.FinishedAt() != nil {
		t.Fatal("antes da aprovacao, nenhum timestamp deve existir")
	}

	_ = so.UpdateStatus(StatusInExecution)
	startedCapture := so.StartedAt()
	if startedCapture == nil {
		t.Fatal("StartedAt() deveria estar gravado apos aprovacao")
	}

	_ = so.UpdateStatus(StatusFinished)
	if so.StartedAt() == nil || !so.StartedAt().Equal(*startedCapture) {
		t.Error("StartedAt() nao deveria ter mudado ao finalizar")
	}
	if so.FinishedAt() == nil {
		t.Error("FinishedAt() deveria estar gravado apos finalizar")
	}
	if so.FinishedAt().Before(*startedCapture) {
		t.Error("FinishedAt() nao pode ser anterior a StartedAt()")
	}
}

// Helper para criar OS com status especifico (via reconstituicao).
func newServiceOrderWithStatus(status OrderStatus) *ServiceOrder {
	t := time.Now()
	return ReconstituteServiceOrder(
		"order-1", 0, "cust-1", "veh-1",
		status,
		nil, nil, 0, "",
		t, t,
		nil, nil,
	)
}
