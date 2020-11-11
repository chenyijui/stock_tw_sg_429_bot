package service

import (
	"context"
	"github.com/chromedp/chromedp"
	TGBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	Define "stock_tw_sg_429_bot/define"
	"time"
)

var res string

func GetInvestmentTrustNetBuy() interface{} {

	url := Define.STOCK_INVESTMENT_TRUST_NET_BUY

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
		//chromedp.EmulateViewport(560, 360, chromedp.EmulateScale(1)),
		chromedp.Navigate(url),
		chromedp.Sleep(1*time.Second),
		//chromedp.CaptureScreenshot(&buf),
		chromedp.Screenshot("div.report-table_wrapper", &buf, chromedp.ByID),
	)
	if cdpErr != nil {
		log.Fatal(cdpErr)
	}

	return TGBotAPI.FileBytes{Name: "image.jpg", Bytes: buf}
}
