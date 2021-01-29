package main

import (
	"fmt"
	botAPI "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shopspring/decimal"
	"log"
	"regexp"
	Common "stock_tw_sg_429_bot/define"
	Define "stock_tw_sg_429_bot/define"
	eStockType "stock_tw_sg_429_bot/enum/e_stock_type"
	Service "stock_tw_sg_429_bot/service"
	Structure "stock_tw_sg_429_bot/structure"
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
			case "/chart_tse":
				{
					handleChart(update, bot, Define.TSE_ID, eStockType.StockChart)
				}
			case "/k_tse":
				{
					handleChart(update, bot, Define.TSE_ID, eStockType.Candlestick)
				}
			}

			continue
		}
		stockNumStr := split[1]
		if isStockIDValid(stockNumStr) {
			msg := botAPI.NewMessage(update.Message.Chat.ID, "股號錯誤")
			bot.Send(msg)
			continue
		}
		switch split[0] {
		case "/chart":
			{
				handleChart(update, bot, stockNumStr, eStockType.StockChart)
			}
		case "/pv":
			{
				s := Service.GetFivePriceAndVolume(stockNumStr)
				fmt.Print(s)
				fmt.Print(s.Mem["name"])
				message := botAPI.NewMessage(update.Message.Chat.ID, "res")
				bot.Send(message)
			}
		case "/k":
			{
				handleChart(update, bot, stockNumStr, eStockType.Candlestick)
			}
		default:
			continue
		}
	}

}

func handleChart(update botAPI.Update, bot *botAPI.BotAPI, stockNumStr string, stockType eStockType.StockType) {
	msgCh := make(chan bool)
	uploadCh := make(chan bool)
	priceAndVolumeCh := make(chan bool)
	errorCh := make(chan bool)
	loadingMsg := botAPI.Message{}
	var sentMsgErr error
	priceAndVolume := Structure.PriceAndVolumeST{}
	go func() {
		msg := botAPI.NewMessage(update.Message.Chat.ID, "loading...")
		loadingMsg, sentMsgErr = bot.Send(msg)
		if sentMsgErr != nil {
			msgCh <- false
		}
		if loadingMsg.Chat.ID != 0 {
			msgCh <- true
		}
	}()
	var context string
	go func() {
		priceAndVolume = Service.GetFivePriceAndVolume(stockNumStr)
		if len(priceAndVolume.Tick) != 0 {
			startPrice := priceAndVolume.Mem["129"].(float64)
			dealPrice := priceAndVolume.Mem["125"].(float64)
			current := "📈"
			preSymbol := "+"
			calNumDec := decimal.NewFromFloat(dealPrice).Sub(decimal.NewFromFloat(startPrice))
			percentDec := decimal.NewFromFloat(dealPrice).Sub(decimal.NewFromFloat(startPrice)).Div(decimal.NewFromFloat(startPrice)).Mul(decimal.New(int64(100), 0))

			if dealPrice < startPrice {
				current = "📉"
				preSymbol = "-"
				percentDec = decimal.NewFromFloat(startPrice).Div(decimal.NewFromFloat(dealPrice))
				calNumDec = decimal.NewFromFloat(startPrice).Sub(decimal.NewFromFloat(dealPrice))
			}
			context = current + " " + priceAndVolume.Mem["id"].(string) + " " + priceAndVolume.Mem["name"].(string) + " " + strconv.FormatFloat(dealPrice, 'f', 2, 64) + " " + preSymbol + tool.DecimalToString(calNumDec) + " (" + preSymbol + tool.DecimalToString(percentDec) + "%)"
			priceAndVolumeCh <- true
		} else {
			if <-msgCh {
				deleteMessage := botAPI.NewDeleteMessage(update.Message.Chat.ID, loadingMsg.MessageID)
				bot.Send(deleteMessage)
				errorCh <- true
			}
		}
	}()
	uploadMessage := botAPI.Message{}
	var uploadMsgErr error
	go func() {
		buff := Service.GetTimeNowStockChart(update, stockNumStr, stockType)
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
	go func() {
		if <-errorCh {
			msg := botAPI.NewMessage(update.Message.Chat.ID, "股號錯誤")
			bot.Send(msg)
		}
	}()
}

func isStockIDValid(stockID string) bool {
	stockIDStr := strings.Replace(stockID, " ", "", -1)
	if len(stockIDStr) == 4 {
		return false
	}
	match, _ := regexp.MatchString("/[0-9][0-9][0-9][0-9]/", stockIDStr)
	return match
}
