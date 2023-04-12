package backtesting

type Tester interface {
	RunTest() (*Report, error)
	Advance()
	AnyLeft() bool
}
