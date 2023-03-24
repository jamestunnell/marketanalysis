package collection

const (
	Resolution1Min = "1m"
)

type Info struct {
	Symbol     string `json:"symbol"`
	Resolution string `json:"resolution"`
}

func NewInfo(sym, res string) *Info {
	return &Info{
		Symbol:     sym,
		Resolution: res,
	}
}
