package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"net"
	"os"
	"runtime"
	"sync"
	"time"

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

var wg sync.WaitGroup

func main() {
	// local message
	//lister - серверный сокет
	listener, err := net.Listen("tcp4", "10.9.38.106:1027")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(listener)

	// ожидаю новые подключения "слушать порт"

	connection, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(connection)
	defer connection.Close()
	t := time.Now().Add(time.Second)
	connection.SetReadDeadline(t)
	msg_out := []byte("Prvet Misha")
	msg_in := make([]byte, 10)

	connection.Write(msg_out)
	n, err := connection.Read(msg_in)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(n)
	fmt.Println(string(msg_in))

	//у клиента не завершилось соедение без отправки пакета на завершение
	file, err := os.Open("token.json")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(runtime.NumCPU())
	decoder := json.NewDecoder(file)
	token := TelegramToken{}
	decoder.Decode(&token)

	wg.Add(2)
	go math(10) //рутина
	go math(20)
	fmt.Println(token.BotToken)

	bot, err := tgbotapi.NewBotAPI(token.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	wg.Wait()
	qrcode.WriteColorFile("https://www.example.com", qrcode.Medium, 256, color.CMYK{200, 20, 60, 10}, color.White, "qr.png")
}

func math(id int) {
	defer wg.Done()
	for i := 1; i < 4; i++ {

		fmt.Printf("%d:%d\n", id, i)
		time.Sleep(time.Second)

	}
}
