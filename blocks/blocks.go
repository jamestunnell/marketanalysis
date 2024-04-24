package blocks

func init() {
	reg.Add(TypeAroon, NewAroon)
	reg.Add(TypeATR, NewATR)
	reg.Add(TypeBar, NewBar)
	reg.Add(TypeDMI, NewDMI)
	reg.Add(TypeEMA, NewEMA)
	reg.Add(TypeEMV, NewEMV)
	reg.Add(TypeHeikinAshi, NewHeikinAshi)
	reg.Add(TypeMAOrder, NewMAOrder)
	reg.Add(TypeSMA, NewSMA)
	reg.Add(TypeSupertrend, NewSupertrend)
}
