package service

import (
	TGBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
	Jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"net/http"
	Common "stock_tw_sg_429_bot/common"
	Structure "stock_tw_sg_429_bot/structure"
)

func GetFivePriceAndVolume(update TGBotAPI.Update, stockNumberStr string) interface{} {

	url := Common.STOCK_PRICE_AND_VOLUME_API_URL + stockNumberStr + "&callback=cb"
	res, getErr := http.Get(url)
	if getErr != nil {
		log.Fatal(getErr)
	}

	bytes, _ := ioutil.ReadAll(res.Body)

	var PriceAndVolumeST Structure.PriceAndVolumeST
	getErr = Jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(bytes, &PriceAndVolumeST)

	return PriceAndVolumeST
}
