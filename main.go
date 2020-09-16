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

		// navigate to a page, wait for an element, click
		var example string
		cdpErr := chromedp.Run(ctx,
			//访问打开必应页面
			chromedp.Navigate(`https://cn.bing.com/?mkt=zh-CN`),
			// 等待右下角图标加载完成
			chromedp.WaitVisible(`#sh_cp_in`),
			//搜索框内输入zhangguanzhang
			chromedp.SendKeys(`#sb_form_q`, `zhangguanzhang`, chromedp.ByID),
			// 点击搜索图标
			chromedp.Click(`#sb_form_go`, chromedp.NodeVisible),
			// 获取第一个搜索结构的超链接
			chromedp.Text(`#b_results > li:nth-child(2) > div > div > cite`, &example),
			chromedp.CaptureScreenshot(&buf),
		)
		if cdpErr != nil {
			log.Fatal(cdpErr)
		}
		if err := ioutil.WriteFile("fullScreenshot.png", buf, 0644); err != nil {
			log.Fatal(err)
		}
		log.Printf("example: %s", example)

		//BotAPI.NewPhotoUpload()
		//msg.ReplyToMessageID = update.Message.MessageID


		upload := BotAPI.NewPhotoUpload(update.Message.Chat.ID, "fullScreenshot.png")


		bot.Send(upload)


	}
}
