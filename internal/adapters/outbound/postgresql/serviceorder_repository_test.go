package postgresql

import (
	"context"
	"errors"
	"testing"
	"time"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
)

func seedVehicle(t *testing.T, plate string) (*entities.Requester, *entities.Vehicle) {
	t.Helper()
	crepo := NewRequesterRepository(testDB)
	doc := "99999999909"
	now := time.Now()
	c := entities.ReconstituteRequester("", "Cliente OS", doc, "os@test.com", "11999990010", now, now)
	if err := crepo.Create(context.Background(), c); err != nil {
		existing, ferr := crepo.FindByDocument(context.Background(), doc)
		if ferr != nil {
			t.Fatalf("seedVehicle: requester create failed: %v", err)
		}
		c = existing
	} else {
		t.Cleanup(func() { testDB.Delete(&pgmodel.Requester{}, "id = ?", c.ID()) })
	}

	vrepo := NewVehicleRepository(testDB)
	v := entities.ReconstituteVehicle("", c.ID(), plate, "Toyota", "Hilux", 2023, now, now)
	if err := vrepo.Create(context.Background(), v); err != nil {
		t.Fatalf("seedVehicle: vehicle create failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Vehicle{}, "id = ?", v.ID()) })

	return c, v
}

func newTestOrder(requesterID, vehicleID string) *entities.ServiceOrder {
	so := entities.NewServiceOrder(requesterID, vehicleID)
	return so
}

func TestServiceOrderRepository_Create(t *testing.T) {
	c, v := seedVehicle(t, "OS00001")
	repo := NewServiceOrderRepository(testDB)
	so := newTestOrder(c.ID(), v.ID())

	if err := repo.Create(context.Background(), so); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if so.ID() == "" {
		t.Fatal("expected ID to be set after create")
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.ServiceOrder{}, "id = ?", so.ID()) })
}

func TestServiceOrderRepository_FindByID(t *testing.T) {
	c, v := seedVehicle(t, "OS00002")
	repo := NewServiceOrderRepository(testDB)
	so := newTestOrder(c.ID(), v.ID())
	if err := repo.Create(context.Background(), so); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.ServiceOrder{}, "id = ?", so.ID()) })

	found, err := repo.FindByID(context.Background(), so.ID())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.ID() != so.ID() {
		t.Errorf("expected ID %s, got %s", so.ID(), found.ID())
	}
}

func TestServiceOrderRepository_FindByID_NotFound(t *testing.T) {
	repo := NewServiceOrderRepository(testDB)

	_, err := repo.FindByID(context.Background(), "00000000-0000-0000-0000-000000000000")
	if !isDomainNotFound(err) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestServiceOrderRepository_FindAll(t *testing.T) {
	c, v := seedVehicle(t, "OS00003")
	repo := NewServiceOrderRepository(testDB)

	soReceived := newTestOrder(c.ID(), v.ID())
	soInDiagnosis := newTestOrder(c.ID(), v.ID())
	soAwaiting := newTestOrder(c.ID(), v.ID())
	soInExecution := newTestOrder(c.ID(), v.ID())
	soFinished := newTestOrder(c.ID(), v.ID())
	soDelivered := newTestOrder(c.ID(), v.ID())

	for _, so := range []*entities.ServiceOrder{soReceived, soInDiagnosis, soAwaiting, soInExecution, soFinished, soDelivered} {
		if err := repo.Create(context.Background(), so); err != nil {
			t.Fatalf("setup Create: %v", err)
		}
		id := so.ID()
		t.Cleanup(func() { testDB.Delete(&pgmodel.ServiceOrder{}, "id = ?", id) })
	}

	testDB.Model(&pgmodel.ServiceOrder{}).Where("id = ?", soInDiagnosis.ID()).Update("status", entities.StatusInDiagnosis)
	testDB.Model(&pgmodel.ServiceOrder{}).Where("id = ?", soAwaiting.ID()).Update("status", entities.StatusAwaitingApproval)
	testDB.Model(&pgmodel.ServiceOrder{}).Where("id = ?", soInExecution.ID()).Update("status", entities.StatusInExecution)
	testDB.Model(&pgmodel.ServiceOrder{}).Where("id = ?", soFinished.ID()).Update("status", entities.StatusFinished)
	testDB.Model(&pgmodel.ServiceOrder{}).Where("id = ?", soDelivered.ID()).Update("status", entities.StatusDelivered)

	all, err := repo.FindAll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, so := range all {
		if so.Status() == entities.StatusFinished || so.Status() == entities.StatusDelivered {
			t.Errorf("FindAll retornou OS com status %s — deveria ser excluida", so.Status())
		}
	}

	indexOf := func(id string) int {
		for i, so := range all {
			if so.ID() == id {
				return i
			}
		}
		return -1
	}

	idxExec := indexOf(soInExecution.ID())
	idxAwait := indexOf(soAwaiting.ID())
	idxDiag := indexOf(soInDiagnosis.ID())
	idxRecv := indexOf(soReceived.ID())

	if idxExec == -1 || idxAwait == -1 || idxDiag == -1 || idxRecv == -1 {
		t.Fatal("uma ou mais OS ativas nao foram retornadas")
	}

	if idxExec >= idxAwait {
		t.Errorf("in_execution (pos %d) deveria vir antes de awaiting_approval (pos %d)", idxExec, idxAwait)
	}
	if idxAwait >= idxDiag {
		t.Errorf("awaiting_approval (pos %d) deveria vir antes de in_diagnosis (pos %d)", idxAwait, idxDiag)
	}
	if idxDiag >= idxRecv {
		t.Errorf("in_diagnosis (pos %d) deveria vir antes de received (pos %d)", idxDiag, idxRecv)
	}
}

func TestServiceOrderRepository_FindByRequesterID(t *testing.T) {
	c, v := seedVehicle(t, "OS00004")
	repo := NewServiceOrderRepository(testDB)
	so := newTestOrder(c.ID(), v.ID())
	if err := repo.Create(context.Background(), so); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.ServiceOrder{}, "id = ?", so.ID()) })

	orders, err := repo.FindByRequesterID(context.Background(), c.ID())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(orders) == 0 {
		t.Fatal("expected at least one service order for this requester")
	}
}

func TestServiceOrderRepository_UpdateStatus(t *testing.T) {
	c, v := seedVehicle(t, "OS00005")
	repo := NewServiceOrderRepository(testDB)
	so := newTestOrder(c.ID(), v.ID())
	if err := repo.Create(context.Background(), so); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.ServiceOrder{}, "id = ?", so.ID()) })

	if err := repo.UpdateStatus(context.Background(), so.ID(), entities.StatusInExecution); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, _ := repo.FindByID(context.Background(), so.ID())
	if found.Status() != entities.StatusInExecution {
		t.Errorf("expected status InProgress, got %s", found.Status())
	}
}

func TestServiceOrderRepository_Update(t *testing.T) {
	c, v := seedVehicle(t, "OS00006")
	repo := NewServiceOrderRepository(testDB)
	so := newTestOrder(c.ID(), v.ID())
	if err := repo.Create(context.Background(), so); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.ServiceOrder{}, "id = ?", so.ID()) })

	so.SetNotes("Observação atualizada")
	if err := repo.Update(context.Background(), so); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, _ := repo.FindByID(context.Background(), so.ID())
	if found.Notes() != "Observação atualizada" {
		t.Errorf("expected updated notes, got %s", found.Notes())
	}
}

func TestServiceOrderRepository_Delete(t *testing.T) {
	c, v := seedVehicle(t, "OS00007")
	repo := NewServiceOrderRepository(testDB)
	so := newTestOrder(c.ID(), v.ID())
	if err := repo.Create(context.Background(), so); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if err := repo.Delete(context.Background(), so.ID()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err := repo.FindByID(context.Background(), so.ID())
	if !isDomainNotFound(err) {
		t.Fatal("expected record to be deleted")
	}
}

func TestServiceOrderRepository_AverageExecutionTime(t *testing.T) {
	c, v := seedVehicle(t, "OS00008")
	repo := NewServiceOrderRepository(testDB)

	now := time.Now().UTC().Truncate(time.Second)
	started1 := now.Add(-2 * time.Hour)
	finished1 := now.Add(-1 * time.Hour) // duracao 1h
	started2 := now.Add(-4 * time.Hour)
	finished2 := now.Add(-1 * time.Hour) // duracao 3h

	so1 := entities.ReconstituteServiceOrder("", 0, c.ID(), v.ID(),
		entities.StatusFinished, nil, nil, 0, "", now, now, &started1, &finished1)
	so2 := entities.ReconstituteServiceOrder("", 0, c.ID(), v.ID(),
		entities.StatusFinished, nil, nil, 0, "", now, now, &started2, &finished2)
	so3 := entities.ReconstituteServiceOrder("", 0, c.ID(), v.ID(),
		entities.StatusReceived, nil, nil, 0, "", now, now, nil, nil)

	for _, so := range []*entities.ServiceOrder{so1, so2, so3} {
		if err := repo.Create(context.Background(), so); err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		t.Cleanup(func() { testDB.Delete(&pgmodel.ServiceOrder{}, "id = ?", so.ID()) })
	}

	avg, err := repo.AverageExecutionTime(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := 2 * time.Hour
	delta := avg - expected
	if delta < 0 {
		delta = -delta
	}
	if delta > time.Second {
		t.Errorf("media = %v, esperava ~%v (delta %v)", avg, expected, delta)
	}
}

func TestServiceOrderRepository_FindByCode(t *testing.T) {
	c, v := seedVehicle(t, "OS00010")
	repo := NewServiceOrderRepository(testDB)
	so := newTestOrder(c.ID(), v.ID())
	if err := repo.Create(context.Background(), so); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.ServiceOrder{}, "id = ?", so.ID()) })

	found, err := repo.FindByCode(context.Background(), so.Code())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.ID() != so.ID() {
		t.Errorf("expected ID %s, got %s", so.ID(), found.ID())
	}
}

func TestServiceOrderRepository_FindByCode_NotFound(t *testing.T) {
	repo := NewServiceOrderRepository(testDB)
	_, err := repo.FindByCode(context.Background(), 999999999)
	if !isDomainNotFound(err) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestServiceOrderRepository_AverageExecutionTime_SemAmostras(t *testing.T) {
	repo := NewServiceOrderRepository(testDB)

	testDB.Exec("DELETE FROM service_orders WHERE started_at IS NOT NULL AND finished_at IS NOT NULL")

	avg, err := repo.AverageExecutionTime(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if avg != 0 {
		t.Errorf("esperava 0 sem amostras, obteve %v", avg)
	}
}

// =============================================================================
// ApplyApprovalTransition — decremento de estoque + save da OS em uma unica
// transacao. Falha em qualquer passo dispara rollback de tudo.
// =============================================================================

// seedPart cria uma peca com o estoque informado e cadastra cleanup.
func seedPart(t *testing.T, code string, stock int) *entities.Part {
	t.Helper()
	repo := NewPartRepository(testDB)
	p := entities.NewPart(code, "Peca de teste "+code, "Desc", "un", 10.00, stock)
	if err := repo.Create(context.Background(), p); err != nil {
		t.Fatalf("seedPart(%s): %v", code, err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.Part{}, "id = ?", p.ID()) })
	return p
}

func partStock(t *testing.T, id string) int {
	t.Helper()
	var m pgmodel.Part
	if err := testDB.First(&m, "id = ?", id).Error; err != nil {
		t.Fatalf("partStock(%s): %v", id, err)
	}
	return m.Stock
}

func TestServiceOrderRepository_ApplyApprovalTransition_Sucesso(t *testing.T) {
	c, v := seedVehicle(t, "OS10001")
	part1 := seedPart(t, "AAT-001", 10)
	part2 := seedPart(t, "AAT-002", 5)

	repo := NewServiceOrderRepository(testDB)
	so := newTestOrder(c.ID(), v.ID())
	so.SetParts([]entities.PartItem{
		{PartID: part1.ID(), Description: part1.Name(), Quantity: 3, UnitPrice: part1.Price()},
		{PartID: part2.ID(), Description: part2.Name(), Quantity: 2, UnitPrice: part2.Price()},
	})
	if err := repo.Create(context.Background(), so); err != nil {
		t.Fatalf("setup Create: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.ServiceOrder{}, "id = ?", so.ID()) })

	if err := so.UpdateStatus(entities.StatusInDiagnosis); err != nil {
		t.Fatalf("transicao diag: %v", err)
	}
	if err := so.UpdateStatus(entities.StatusAwaitingApproval); err != nil {
		t.Fatalf("transicao awaiting: %v", err)
	}
	if err := so.UpdateStatus(entities.StatusInExecution); err != nil {
		t.Fatalf("transicao exec: %v", err)
	}

	if err := repo.ApplyApprovalTransition(context.Background(), so); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got := partStock(t, part1.ID()); got != 7 {
		t.Errorf("part1 stock = %d, esperava 7 (10 - 3)", got)
	}
	if got := partStock(t, part2.ID()); got != 3 {
		t.Errorf("part2 stock = %d, esperava 3 (5 - 2)", got)
	}
	persisted, err := repo.FindByID(context.Background(), so.ID())
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if persisted.Status() != entities.StatusInExecution {
		t.Errorf("status persistido = %q, esperava in_execution", persisted.Status())
	}
	if persisted.StartedAt() == nil {
		t.Error("StartedAt nao foi persistido")
	}
}

func TestServiceOrderRepository_ApplyApprovalTransition_EstoqueInsuficiente_RollbackTotal(t *testing.T) {
	c, v := seedVehicle(t, "OS10002")
	part1 := seedPart(t, "AAT-011", 10) // suficiente
	part2 := seedPart(t, "AAT-012", 1)  // insuficiente (pede 5)

	repo := NewServiceOrderRepository(testDB)
	so := newTestOrder(c.ID(), v.ID())
	so.SetParts([]entities.PartItem{
		{PartID: part1.ID(), Description: part1.Name(), Quantity: 3, UnitPrice: part1.Price()},
		{PartID: part2.ID(), Description: part2.Name(), Quantity: 5, UnitPrice: part2.Price()},
	})
	if err := repo.Create(context.Background(), so); err != nil {
		t.Fatalf("setup Create: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.ServiceOrder{}, "id = ?", so.ID()) })

	if err := so.UpdateStatus(entities.StatusInDiagnosis); err != nil {
		t.Fatalf("transicao diag: %v", err)
	}
	if err := so.UpdateStatus(entities.StatusAwaitingApproval); err != nil {
		t.Fatalf("transicao awaiting: %v", err)
	}
	if err := so.UpdateStatus(entities.StatusInExecution); err != nil {
		t.Fatalf("transicao exec: %v", err)
	}

	err := repo.ApplyApprovalTransition(context.Background(), so)
	if !errors.Is(err, domainerrors.ErrInsufficientStock) {
		t.Fatalf("erro = %v, esperava ErrInsufficientStock", err)
	}

	if got := partStock(t, part1.ID()); got != 10 {
		t.Errorf("part1 stock = %d, esperava 10 (rollback)", got)
	}
	if got := partStock(t, part2.ID()); got != 1 {
		t.Errorf("part2 stock = %d, esperava 1 (inalterado)", got)
	}
	persisted, ferr := repo.FindByID(context.Background(), so.ID())
	if ferr != nil {
		t.Fatalf("FindByID: %v", ferr)
	}
	if persisted.Status() == entities.StatusInExecution {
		t.Error("status no banco mudou para in_execution apesar do erro — rollback falhou")
	}
}

func TestServiceOrderRepository_ApplyApprovalTransition_PecaInexistente_Rollback(t *testing.T) {
	c, v := seedVehicle(t, "OS10003")
	part1 := seedPart(t, "AAT-021", 10)

	repo := NewServiceOrderRepository(testDB)
	so := newTestOrder(c.ID(), v.ID())
	so.SetParts([]entities.PartItem{
		{PartID: part1.ID(), Description: part1.Name(), Quantity: 2, UnitPrice: part1.Price()},
		{PartID: "00000000-0000-0000-0000-000000000000", Description: "inexistente", Quantity: 1, UnitPrice: 10.00},
	})
	if err := repo.Create(context.Background(), so); err != nil {
		t.Fatalf("setup Create: %v", err)
	}
	t.Cleanup(func() { testDB.Delete(&pgmodel.ServiceOrder{}, "id = ?", so.ID()) })

	if err := so.UpdateStatus(entities.StatusInDiagnosis); err != nil {
		t.Fatalf("transicao diag: %v", err)
	}
	if err := so.UpdateStatus(entities.StatusAwaitingApproval); err != nil {
		t.Fatalf("transicao awaiting: %v", err)
	}
	if err := so.UpdateStatus(entities.StatusInExecution); err != nil {
		t.Fatalf("transicao exec: %v", err)
	}

	err := repo.ApplyApprovalTransition(context.Background(), so)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("erro = %v, esperava ErrNotFound", err)
	}
	if got := partStock(t, part1.ID()); got != 10 {
		t.Errorf("part1 stock = %d, esperava 10 (rollback da peca que existia)", got)
	}
}
