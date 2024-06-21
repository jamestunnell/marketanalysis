package background

type Job interface {
	GetID() string
	Execute(progress JobProgressFunc) (Cloneable, error)
}

type JobProgressFunc func(update ProgressUpdate)
