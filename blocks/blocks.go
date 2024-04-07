package blocks

func init() {
	reg.Add(TypeAroon, NewAroon)
	reg.Add(TypeATR, NewATR)
	reg.Add(TypeDMI, NewDMI)
	reg.Add(TypeEMA, NewEMA)
	reg.Add(TypeEMV, NewEMV)
	reg.Add(TypeSMA, NewSMA)
}
