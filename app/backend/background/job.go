package background

type Job interface {
	GetID() string
	Execute(progress JobProgressFunc) (any, error)
}

type JobProgressFunc func(update ProgressUpdate)
