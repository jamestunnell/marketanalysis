package optimization

type Settings struct {
	Algorithm     string `json:"algorithm"`
	MaxIterations int    `json:"maxIterations"`
	KeepHistory   bool   `json:"keepHistory"`
}
