package background

type Job interface {
	GetID() string
	Execute(onProgress JobProgressFunc) (any, error)
}

type JobProgressFunc func(progress float64)
