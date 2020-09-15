package main

import (
	BotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	Common "stock_tw_sg_429_bot/common"
)

func main() {
	bot, err := BotAPI.NewBotAPI(Common.BOT_TOKEN)
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

		msg := BotAPI.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}