package service

import (
	Jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"net/http"
	Common "stock_tw_sg_429_bot/define"
	Structure "stock_tw_sg_429_bot/structure"
)

func GetFivePriceAndVolume(stockNumberStr string) Structure.PriceAndVolumeST {

	url := Common.STOCK_PRICE_AND_VOLUME_API_URL + stockNumberStr + "&callback=cb"
	res, getErr := http.Get(url)
	if getErr != nil {
		log.Fatal(getErr)
	}

	bytes, _ := ioutil.ReadAll(res.Body)
	i := len(bytes)
	i2 := bytes[3 : i-1]

	var PriceAndVolumeST Structure.PriceAndVolumeST
	getErr = Jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(i2, &PriceAndVolumeST)

	return PriceAndVolumeST
}
