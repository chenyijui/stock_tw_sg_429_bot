package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	BotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	Common "stock_tw_sg_429_bot/common"
	"strings"
	"time"
)

func main() {
	bot, err := BotAPI.NewBotAPI(Common.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := BotAPI.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		split := strings.Split(update.Message.Text, "/chart")
		fmt.Print(split[1])
		stockNum := strings.Replace(split[1], " ", "", -1)
		API := "https://s.yimg.com/nb/tw_stock_frontend/scripts/StxChart/StxChart.9d11dfe155.html?sid=" + stockNum
		res, err := http.Get(API)
		if err != nil {
			log.Fatal(err)
		}

		s, _ := ioutil.ReadAll(res.Body)

		fmt.Print(string(s))
		update.Message.Text = string(s)
		//msg := BotAPI.NewMessage(update.Message.Chat.ID, update.Message.Text)

		// start get photo

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
			chromedp.Navigate(`https://s.yimg.com/nb/tw_stock_frontend/scripts/StxChart/StxChart.9d11dfe155.html?sid=2330`),
			chromedp.CaptureScreenshot(&buf),
		)
		if cdpErr != nil {
			log.Fatal(cdpErr)
		}

		b := BotAPI.FileBytes{Name: "image.jpg", Bytes: buf}


		//BotAPI.NewPhotoUpload()
		//msg.ReplyToMessageID = update.Message.MessageID

		upload := BotAPI.NewPhotoUpload(update.Message.Chat.ID, b)


		bot.Send(upload)

	}
}
