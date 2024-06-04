package backend

import "github.com/jamestunnell/marketanalysis/models"

type Synchronizer interface {
	Start()
	Stop()
	Trigger(trig *SyncTrigger)
}

type SyncTrigger struct {
	Security *models.Security
	Op       SyncOp
}

type SyncOp int

const (
	SyncAdd SyncOp = iota
	SyncRemove
	SyncScan
)

func TriggerAdd(s *models.Security) *SyncTrigger {
	return &SyncTrigger{Security: s, Op: SyncAdd}
}

func TriggerRemove(sym string) *SyncTrigger {
	s := &models.Security{Symbol: sym}

	return &SyncTrigger{Security: s, Op: SyncRemove}
}

func TriggerScan(s *models.Security) *SyncTrigger {
	return &SyncTrigger{Security: s, Op: SyncScan}
}

func (op SyncOp) String() string {
	var s string

	switch op {
	case SyncAdd:
		s = "add"
	case SyncRemove:
		s = "remove"
	case SyncScan:
		s = "scan"
	}

	return s
}
