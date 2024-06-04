package api

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"

	"github.com/jamestunnell/marketanalysis/app/backend"
	"github.com/jamestunnell/marketanalysis/app/backend/stores"
	"github.com/jamestunnell/marketanalysis/models"
)

type Securities struct {
	*CRUDAPI[*models.Security]
	DB   *mongo.Database
	Sync backend.Synchronizer
}

func NewSecurities(
	db *mongo.Database,
	synchronizer backend.Synchronizer,
) *Securities {
	s := &Securities{
		DB:   db,
		Sync: synchronizer,
	}

	mods := []CRUDOptionsMod[*models.Security]{
		WithPostCreateHook(s.triggerSyncAdd),
		WithPostUpdateHook(s.triggerSyncScan),
		WithPostDeleteHook[*models.Security](s.triggerSyncRemove),
	}

	s.CRUDAPI = NewCRUDAPI(stores.NewSecurities(db), mods...)

	return s
}

func (a *Securities) triggerSyncAdd(sec *models.Security) {
	a.Sync.Trigger(backend.TriggerAdd(sec))
}

func (a *Securities) triggerSyncScan(sec *models.Security) {
	a.Sync.Trigger(backend.TriggerScan(sec))
}

func (a *Securities) triggerSyncRemove(symbol string) {
	a.Sync.Trigger(backend.TriggerRemove(symbol))
}

func (a *Securities) Bind(r *mux.Router) {
	a.CRUDAPI.Bind(r)

	// r.HandleFunc(a.SingularRoute()+"/collect", a.BacktestGraph).Methods(http.MethodPost) //, http.MethodOptions)
	// r.HandleFunc(a.SingularRoute()+"/eval", a.EvalGraph).Methods(http.MethodPost)        //, http.MethodOptions)
	// r.HandleFunc(a.SingularRoute()+"/run", a.RunGraph).Methods(http.MethodPost)          //, http.MethodOptions)
}
