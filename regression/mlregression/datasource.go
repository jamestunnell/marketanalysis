package mlregression

type Data interface {
	Inputs() [][]float64
	Output() []float64
}
