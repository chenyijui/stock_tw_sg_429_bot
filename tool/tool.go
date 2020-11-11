package tool

import (
	"github.com/shopspring/decimal"
	"strconv"
)

func DecimalToString(decimal decimal.Decimal) string {
	num, _ := decimal.Float64()
	return strconv.FormatFloat(num, 'f', 2, 64)
}
