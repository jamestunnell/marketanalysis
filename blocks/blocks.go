package blocks

func init() {
	registry.Add(TypeAroon, NewAroon)
	registry.Add(TypeATR, NewATR)
	registry.Add(TypeDMI, NewDMI)
	registry.Add(TypeEMA, NewEMA)
	registry.Add(TypeEMV, NewEMV)
	registry.Add(TypeSMA, NewSMA)
}
