package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/onionj/pricebot/price"
	"github.com/onionj/pricebot/telegram"
)

const (
	NEW_MESSAGE_PERIOD    = 60 * 60 * 4
	UPDATE_MESSAGE_PERIOD = 5
	UPDATE_PRICE_PERIOD   = 60
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	BOT_TOKEN := os.Getenv("BOT_TOKEN")
	CHAT_ID := os.Getenv("CHAT_ID")
	CHANEL_NAME := os.Getenv("CHANEL_NAME")
	PROXY_LINK := os.Getenv("PROXY_LINK")

	if BOT_TOKEN == "" || CHAT_ID == "" {
		fmt.Println("Missing BOT_TOKEN or CHAT_ID in environment variables")
		return
	}

	price := price.NewPrice()
	tel := telegram.NewTelegram(BOT_TOKEN, CHAT_ID)

	for ; ; time.Sleep(time.Second * UPDATE_MESSAGE_PERIOD) {
		if (time.Now().Unix() - price.LastRefresh.Unix()) >= UPDATE_PRICE_PERIOD {
			err := price.Refresh()
			fmt.Println(price.String())
			if err != nil {
				fmt.Println("refresh price error:", err.Error())
				time.Sleep(time.Minute)
				continue
			}
		}

		nextUpdateSecond := int64(
			math.Min(
				float64(UPDATE_PRICE_PERIOD-(time.Now().Unix()-price.LastRefresh.Unix())),
				UPDATE_PRICE_PERIOD))

		message := createTelegramMessage(price.String(), nextUpdateSecond, CHANEL_NAME, false, PROXY_LINK)

		if tel.LastMessageTime > 0 && tel.LastMessageId > 0 && (time.Now().Unix()-tel.LastMessageTime) <= NEW_MESSAGE_PERIOD {
			if err := tel.UpdateMessage(message, tel.LastMessageId); err != nil {
				fmt.Println("update telegram error:", err.Error())
				time.Sleep(time.Minute)
				continue
			}
		} else {
			tel.UpdateMessage(createTelegramMessage(price.String(), nextUpdateSecond, CHANEL_NAME, true, PROXY_LINK), tel.LastMessageId)

			if err := tel.SendMessage(message); err != nil {
				fmt.Println("send telegram error:", err.Error())
				time.Sleep(time.Minute)
				continue
			}
		}
	}
}

func createTelegramMessage(priceData string, nextUpdateSecond int64, chanelName string, ending bool, proxyLink string) string {
	proxy := fmt.Sprintf("ا🗝 [پروکسی](%s)", proxyLink)

	if ending {
		return fmt.Sprintf("%s\n\n%s", priceData, chanelName)
	}

	if nextUpdateSecond >= 7 {
		return fmt.Sprintf("ا⏰ %02d ثانیه تا بروزرسانی بعدی قیمت ها\n%s\n\n%s\n%s", nextUpdateSecond, priceData, proxy, chanelName)
	} else if nextUpdateSecond >= 3 {
		return fmt.Sprintf("ا🔄 درحال بروزرسانی قیمت ها \n%s\n\n%s\n%s", priceData, proxy, chanelName)
	} else {
		return fmt.Sprintf("ا🔄 درحال بروزرسانی قیمت ها\n%s\n\n%s\n%s", priceData, proxy, chanelName)
	}
}
