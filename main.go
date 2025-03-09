package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/onionj/pricebot/price"
	"github.com/onionj/pricebot/telegram"
)

const (
	NEW_MESSAGE_PERIOD    = 60 * 60 // 1H
	UPDATE_MESSAGE_PERIOD = time.Minute
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	BOT_TOKEN := os.Getenv("BOT_TOKEN")
	CHAT_ID := os.Getenv("CHAT_ID")

	if BOT_TOKEN == "" || CHAT_ID == "" {
		fmt.Println("Missing BOT_TOKEN or CHAT_ID in environment variables")
		return
	}

	price := price.NewPrice()
	tel := telegram.NewTelegram(BOT_TOKEN, CHAT_ID)

	for ; ; time.Sleep(UPDATE_MESSAGE_PERIOD) {
		err := price.Refresh()
		if err != nil {
			fmt.Println("refresh price error:", err.Error())
			time.Sleep(UPDATE_MESSAGE_PERIOD)
			continue
		}
		fmt.Println(price.String(), tel.LastMessageId, tel.LastMessageTime)

		if tel.LastMessageTime > 0 && tel.LastMessageId > 0 && (time.Now().Unix()-tel.LastMessageTime) <= NEW_MESSAGE_PERIOD {
			err = tel.UpdateMessage(price.String(), tel.LastMessageId)
		} else {
			err = tel.SendMessage(price.String())
		}

		if err != nil {
			fmt.Println("send telegram error:", err.Error())
			time.Sleep(UPDATE_MESSAGE_PERIOD)
			continue
		}
	}

}
