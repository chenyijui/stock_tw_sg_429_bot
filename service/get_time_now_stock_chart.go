package service

import (
	"context"
	"github.com/chromedp/chromedp"
	TGBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	Common "stock_tw_sg_429_bot/common"
	"time"
)

func GetTimeNowStockChart(update TGBotAPI.Update, stockNumberStr string) interface{} {

	url := Common.STOCK_CHART_API_URL + stockNumberStr
	res, getErr := http.Get(url)
	if getErr != nil {
		log.Fatal(getErr)
	}

	bodyStr, _ := ioutil.ReadAll(res.Body)

	update.Message.Text = string(bodyStr)

	var buf []byte

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	cdpErr := chromedp.Run(ctx,
		chromedp.EmulateViewport(560, 360, chromedp.EmulateScale(1)),
		chromedp.Navigate(url),
		chromedp.CaptureScreenshot(&buf),
	)
	if cdpErr != nil {
		log.Fatal(cdpErr)
	}

	return TGBotAPI.FileBytes{Name: "image.jpg", Bytes: buf}
}
