package backtesting

type Tester interface {
	RunTest() error
	Advance()
	AnyLeft() bool
}
