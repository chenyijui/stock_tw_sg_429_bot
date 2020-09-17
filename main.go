package main

import (
	"fmt"
	TGBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	Common "stock_tw_sg_429_bot/common"
	Service "stock_tw_sg_429_bot/service"
	"strings"
)

func main() {
	bot, err := TGBotAPI.NewBotAPI(Common.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := TGBotAPI.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		split := strings.Split(update.Message.Text, " ")

		if len(split) != 2 {
			continue
		}
		stockNumStr := strings.Replace(split[1], " ", "", -1)
		switch split[0] {
		case "/chart":
			{
				buff := Service.GetTimeNowStockChart(update, stockNumStr)
				upload := TGBotAPI.NewPhotoUpload(update.Message.Chat.ID, buff)
				bot.Send(upload)
			}
		case "/volum":
			{
				s := Service.GetFivePriceAndVolume(update, stockNumStr)
				fmt.Print(s)
				message := TGBotAPI.NewMessage(update.Message.Chat.ID, "res")
				bot.Send(message)
			}
		default:
			continue
		}

	}
}
