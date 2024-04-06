package blocks

func init() {
	registry.Add(TypeAroon, NewAroon)
	registry.Add(TypeDMI, NewDMI)
	registry.Add(TypeEMA, NewEMA)
	registry.Add(TypeSMA, NewSMA)
}
