package backtesting

type Backtester interface {
	RunTest() (*Report, error)
	Advance()
	AnyLeft() bool
}
