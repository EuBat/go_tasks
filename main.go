package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skip2/go-qrcode"
)

type TelegramToken struct {
	BotToken string
}

func Pic(dx, dy int) [][]uint8 {
	img := make([][]uint8, dy)
	for y := range img {
		img[y] = make([]uint8, dx)
		for x := range img[y] {
			img[y][x] = uint8(x*3 + y*2)
		}
	}
	return img
}

func main() {
	file, err := os.Open("token.json")
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(file)
	token := TelegramToken{}
	decoder.Decode(&token)

	fmt.Println("\nToken:", token.BotToken)

	bot, err := tgbotapi.NewBotAPI(token.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("\nAuthorized on account %s", bot.Self.UserName)
	qrcode.WriteColorFile("https://www.example.com", qrcode.Medium, 256, color.CMYK{200, 20, 60, 10}, color.White, "qr.png")
}
