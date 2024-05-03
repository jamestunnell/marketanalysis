package graph

type Connection struct {
	Source *Address `json:"source"`
	Target *Address `json:"target"`
}
