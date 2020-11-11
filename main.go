package main

import (
	"fmt"
	botAPI "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shopspring/decimal"
	"log"
	Common "stock_tw_sg_429_bot/define"
	Service "stock_tw_sg_429_bot/service"
	"stock_tw_sg_429_bot/tool"
	"strconv"
	"strings"
)

func main() {
	bot, err := botAPI.NewBotAPI(Common.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := botAPI.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		split := strings.Split(update.Message.Text, " ")

		if len(split) != 2 {
			switch split[0] {
			case "/test":
				{
					buff := Service.GetInvestmentTrustNetBuy()
					upload := botAPI.NewPhotoUpload(update.Message.Chat.ID, buff)
					t, _ := bot.Send(upload)
					fmt.Println(t)
				}
			}

			continue
		}
		stockNumStr := strings.Replace(split[1], " ", "", -1)
		switch split[0] {
		case "/chart":
			{
				msgCh := make(chan bool)
				uploadCh := make(chan bool)
				priceAndVolumeCh := make(chan bool)
				loadingMsg := botAPI.Message{}
				var sentMsgErr error

				go func() {
					msg := botAPI.NewMessage(update.Message.Chat.ID, "loading...")
					loadingMsg, sentMsgErr = bot.Send(msg)
					if err != nil {
						msgCh <- false
					}
					if loadingMsg.Chat.ID != 0 {
						msgCh <- true
					}
				}()

				var context string
				go func() {
					priceAndVolume := Service.GetFivePriceAndVolume(stockNumStr)
					startPrice := priceAndVolume.Mem["129"].(float64)
					dealPrice := priceAndVolume.Mem["127"].(float64)
					current := "📈"
					preSymbol := "+"
					calNumDec := decimal.NewFromFloat(dealPrice).Sub(decimal.NewFromFloat(startPrice))
					percentDec := decimal.NewFromFloat(dealPrice).Sub(decimal.NewFromFloat(startPrice)).Div(decimal.NewFromFloat(startPrice)).Mul(decimal.NewFromInt(100))

					if dealPrice < startPrice {
						current = "📉"
						preSymbol = "-"
						percentDec = decimal.NewFromFloat(startPrice).Div(decimal.NewFromFloat(dealPrice))
						calNumDec = decimal.NewFromFloat(startPrice).Sub(decimal.NewFromFloat(dealPrice))
					}
					context = current + " " + priceAndVolume.Mem["id"].(string) + " " + priceAndVolume.Mem["name"].(string) + " " + strconv.FormatFloat(dealPrice, 'f', 2, 64) + " " + preSymbol + tool.DecimalToString(calNumDec) + " (" + preSymbol + tool.DecimalToString(percentDec) + "%)"
					priceAndVolumeCh <- true
				}()

				uploadMessage := botAPI.Message{}
				var uploadMsgErr error
				go func() {

					buff := Service.GetTimeNowStockChart(update, stockNumStr)
					upload := botAPI.NewPhotoUpload(update.Message.Chat.ID, buff)
					if <-priceAndVolumeCh {
						upload.Caption = context
						if <-msgCh {
							deleteMessage := botAPI.NewDeleteMessage(update.Message.Chat.ID, loadingMsg.MessageID)
							bot.Send(deleteMessage)
							uploadMessage, uploadMsgErr = bot.Send(upload)
							if uploadMsgErr != nil {
								uploadCh <- false
							}
							if uploadMessage.Chat.ID != 0 {
								uploadCh <- true
							}
						}
					}
				}()
			}
		case "/pv":
			{
				s := Service.GetFivePriceAndVolume(stockNumStr)
				fmt.Print(s)
				fmt.Print(s.Mem["name"])
				message := botAPI.NewMessage(update.Message.Chat.ID, "res")
				bot.Send(message)
			}
		default:
			continue
		}

	}
}
