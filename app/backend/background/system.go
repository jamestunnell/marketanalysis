package background

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type System interface {
	Start()
	Stop()

	RunJob(Job) bool
	GetJobStatus(id string) (Status, bool)

	Subscribe(sub Subscriber)
	Unsubscribe(subID string)
}

type CheckJobFunc func(Status)

type Subscriber interface {
	GetID() string
	OnUpdate(Status)
}

type system struct {
	stop        chan struct{}
	updates     chan ProgressUpdate
	statusMutex sync.RWMutex
	subscrMutex sync.RWMutex
	status      map[string]Status
	pruneAge    time.Duration
	subscribers map[string]Subscriber
}

type ProgressUpdate struct {
	JobID    string
	Progress float64
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
		status:      map[string]Status{},
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

func (sys *system) RunJob(j Job) bool {
	sys.statusMutex.Lock()

	defer sys.statusMutex.Unlock()

	id := j.GetID()

	_, found := sys.status[id]
	if found {
		return false
	}

	sys.status[id] = Status{
		State:   Running,
		Started: time.Now(),
	}

	go sys.runJob(j)

	return true
}

func (sys *system) GetJobStatus(id string) (Status, bool) {
	sys.statusMutex.RLocker()

	defer sys.statusMutex.RUnlock()

	s, found := sys.status[id]
	if !found {
		return Status{}, false
	}

	return s, true
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
			if sys.updateProgress(update) {
				sys.notifySubs(update.JobID)
			}

		case <-time.After(5 * time.Minute):
			sys.prune()
		}
	}

	log.Info().Msg("background: stopped system")
}

func (sys *system) runJob(j Job) {
	id := j.GetID()

	result, err := j.Execute(func(progress float64) {
		sys.updates <- ProgressUpdate{JobID: id, Progress: progress}
	})

	completed := time.Now()

	sys.statusMutex.Lock()

	defer sys.statusMutex.Unlock()

	s, found := sys.status[id]
	if !found {
		log.Error().Err(err).Msg("background: failed to update status after job completion, job status not found")

		return
	}

	s.Completed = completed

	if err != nil {
		s.State = Failed
		s.ErrMsg = err.Error()
		s.Result = nil
	} else {
		s.State = Succeeded
		s.Result = result
		s.ErrMsg = ""
	}

	sys.status[id] = s
}

func (sys *system) updateProgress(upd ProgressUpdate) bool {
	sys.statusMutex.Lock()

	defer sys.statusMutex.Unlock()

	s, found := sys.status[upd.JobID]
	if !found {
		log.Warn().Str("id", upd.JobID).Msg("background: failed to update progress, job status not found")

		return false
	}

	s.Progress = upd.Progress

	sys.status[upd.JobID] = s

	return true
}

func (sys *system) notifySubs(jobID string) {
	sys.statusMutex.RLock()

	defer sys.statusMutex.RUnlock()

	s, found := sys.status[jobID]
	if !found {
		log.Warn().Str("id", jobID).Msg("background: failed to notify subscribers progress, job status not found")

		return
	}

	for _, sub := range sys.subscribers {
		sub.OnUpdate(s)
	}
}

func (sys *system) prune() {
	sys.statusMutex.Lock()

	defer sys.statusMutex.Unlock()

	toPrune := []string{}

	for id, status := range sys.status {
		if status.State == Running {
			continue
		}

		if time.Since(status.Completed) >= sys.pruneAge {
			toPrune = append(toPrune, id)
		}
	}

	for _, id := range toPrune {
		log.Debug().Str("id", id).Msg("background: pruning job status")

		delete(sys.status, id)
	}
}
