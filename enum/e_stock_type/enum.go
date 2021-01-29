package e_stock_type

// StockType Enum
type StockType int

// StockType Enum編號
const (
	StockChart  StockType = 0 // 即時走勢
	Candlestick StockType = 1 // 蠟燭日K
)

// string fmt
func (code StockType) String() string {
	return string(code)
}

// Int fmt
func (code StockType) Int() int {
	return int(code)
}
