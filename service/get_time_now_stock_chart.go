package service

import (
	"context"
	"github.com/chromedp/chromedp"
	TGBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	Define "stock_tw_sg_429_bot/define"
	eStockType "stock_tw_sg_429_bot/enum/e_stock_type"
	"time"
)

func GetTimeNowStockChart(update TGBotAPI.Update, stockNumberStr string, stockType eStockType.StockType) interface{} {
	var url string
	var height int64
	switch stockType {
	case eStockType.StockChart:
		{
			url = Define.STOCK_CHART_API_URL + stockNumberStr
			height = 360
		}
	case eStockType.Candlestick:
		{
			url = Define.STOCK_CANDLESTICK_API_URL + stockNumberStr
			height = 410
		}
	}

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
		chromedp.EmulateViewport(560, height, chromedp.EmulateScale(1)),
		chromedp.Navigate(url),
		chromedp.CaptureScreenshot(&buf),
	)
	if cdpErr != nil {
		log.Fatal(cdpErr)
	}

	return TGBotAPI.FileBytes{Name: "image.jpg", Bytes: buf}
}
