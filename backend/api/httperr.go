package api

type HTTPErr struct {
	Error      error
	StatusCode int
}
