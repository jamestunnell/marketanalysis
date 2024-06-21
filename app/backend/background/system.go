package background

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type System interface {
	Start()
	Stop()

	Run(Job) error
	GetStatus(id string) (*Status, bool)

	Subscribe(sub Subscriber)
	Unsubscribe(subID string)
}

type CheckJobFunc func(Status)

type Subscriber interface {
	GetID() string
	OnUpdate(*Status)
}

type system struct {
	stop        chan struct{}
	updates     chan ProgressUpdate
	statusMutex sync.RWMutex
	subscrMutex sync.RWMutex
	status      map[string]*Status
	pruneAge    time.Duration
	subscribers map[string]Subscriber
}

type ProgressUpdate struct {
	JobID         string
	InnerProgress float64
	OuterProgress float64
}

type SystemOpts struct {
	PruneAge time.Duration
}

type SystemOptMod func(*SystemOpts)

const DefaultPruneAge = time.Hour

func WithPruneAge(dur time.Duration) SystemOptMod {
	return func(opts *SystemOpts) {
		opts.PruneAge = dur
	}
}

func NewSystem(mods ...SystemOptMod) System {
	opts := &SystemOpts{
		PruneAge: DefaultPruneAge,
	}

	for _, mod := range mods {
		mod(opts)
	}

	return &system{
		stop:        make(chan struct{}),
		status:      map[string]*Status{},
		updates:     make(chan ProgressUpdate),
		pruneAge:    opts.PruneAge,
		subscribers: map[string]Subscriber{},
	}
}

func (sys *system) Start() {
	go sys.runUntilStopped()
}

func (sys *system) Stop() {
	sys.stop <- struct{}{}
}

func (sys *system) Run(j Job) error {
	sys.statusMutex.Lock()

	defer sys.statusMutex.Unlock()

	id := j.GetID()

	_, found := sys.status[id]
	if found {
		return fmt.Errorf("job with ID %s already exists", id)
	}

	sys.status[id] = &Status{State: Queued}

	go sys.runJob(j)

	return nil
}

func (sys *system) GetStatus(id string) (*Status, bool) {
	sys.statusMutex.RLocker()

	defer sys.statusMutex.RUnlock()

	s, found := sys.status[id]
	if !found {
		return nil, false
	}

	return s.Clone(), true
}

func (sys *system) Subscribe(sub Subscriber) {
	sys.subscrMutex.Lock()

	defer sys.subscrMutex.Unlock()

	sys.subscribers[sub.GetID()] = sub
}

func (sys *system) Unsubscribe(subID string) {
	sys.subscrMutex.Lock()

	defer sys.subscrMutex.Unlock()

	delete(sys.subscribers, subID)
}

func (sys *system) runUntilStopped() {
	log.Info().Msg("background: started system")

	keepGoing := true

	for keepGoing {
		select {
		case <-sys.stop:
			keepGoing = false
		case update := <-sys.updates:
			if status := sys.updateProgress(update); status != nil {
				sys.notifySubs(update.JobID, status)
			}
		case <-time.After(5 * time.Minute):
			sys.prune()
		}
	}

	log.Info().Msg("background: stopped system")
}

func (sys *system) runJob(j Job) {
	id := j.GetID()

	result, err := j.Execute(func(update ProgressUpdate) {
		sys.updates <- update
	})

	sys.statusMutex.Lock()

	defer sys.statusMutex.Unlock()

	s, found := sys.status[id]
	if !found {
		log.Error().Err(err).Msg("failed to find job status to update after job completion")

		return
	}

	if err != nil {
		s.ErrMsg = err.Error()
		s.Result = nil
	} else {
		s.Result = result
		s.ErrMsg = ""
	}
}

func (sys *system) updateProgress(upd ProgressUpdate) *Status {
	sys.statusMutex.Lock()

	defer sys.statusMutex.Unlock()

	s, found := sys.status[upd.JobID]
	if !found {
		log.Warn().Str("id", upd.JobID).Msg("background: failed to update progress, job status not found")

		return nil
	}

	s.InnerProgress = upd.InnerProgress
	s.OuterProgress = upd.OuterProgress

	return s.Clone()
}

func (sys *system) notifySubs(jobID string, status *Status) {
	for _, sub := range sys.subscribers {
		sub.OnUpdate(status)
	}
}

func (sys *system) prune() {
	log.Debug().Msg("background: pruning")
}
