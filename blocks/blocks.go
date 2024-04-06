package blocks

func init() {
	registry.Add(TypeSMA, NewSMA)
	registry.Add(TypeEMA, NewEMA)
	registry.Add(TypeAroon, NewAroon)
}
