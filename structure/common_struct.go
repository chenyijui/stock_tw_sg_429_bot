package structure

type PriceAndVolumeST struct {
	Mkt  string `json:"mkt"`
	ID   string `json:"id"`
	Perd string `json:"perd"`
	Type string `json:"type"`
	Mem  string `json:"mem"`
	Tick []Tick `json:"tick"`
}

type Tick struct {
	T float64 `json:"t"`
	P float64 `json:"p"`
	V float64 `json:"v"`
}
