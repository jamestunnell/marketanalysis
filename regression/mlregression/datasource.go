package linregression

type DataSource interface {
	GetData() ([][]float64, []float64, error)
}
