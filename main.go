package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type client struct {
	connection net.Conn
	chanel_in  chan string // input msg
	id         int
}

var (
	clientCount    int = 0
	chanelEntering     = make(chan client)
	chanelExit         = make(chan client)
	chanelMessage      = make(chan string)
)

// печать входящих сообщений
func clientWriter(c client) {
	for {
		c.connection.Write([]byte(<-c.chanel_in + "\n"))
	}
}

// чтение из консоли исходящих сообщений
func clientReader(c client) {
	msg := bufio.NewScanner(c.connection)
	for {
		msg.Scan()
		if msg.Text() == "exit" {
			chanelExit <- c
		} else {
			chanelMessage <- msg.Text()
		}
	}
}

// широковещатель
func broadcaster() {
	clients := make(map[int]client)
	for {
		select {
		case user := <-chanelEntering:
			//add new user to slice
			clients[user.id] = user

			//meeting with the new user
			clients[user.id].chanel_in <- "<- hi #" + strconv.Itoa((clients[user.id].id))
			// tell everybody about new user
			for i := range clients {
				clients[i].chanel_in <- "<- new user #" + strconv.Itoa((clients[user.id].id)) + " connected"
			}
		case user_msg := <-chanelMessage:
			msg := strings.Split(user_msg, " ")
			userId, _ := strconv.Atoi(msg[0])
			for i := range clients {
				if clients[i].id == userId {
					for j := 1; j < len(msg); j++ {
						clients[i].chanel_in <- msg[j]
					}
				}
			}
		case user_exit := <-chanelExit:
			clients[user_exit.id].chanel_in <- "<- bye-bye # " + strconv.Itoa((clients[user_exit.id].id))
			fmt.Println("\n <-client #" + strconv.Itoa((clients[user_exit.id].id)) + " walked")
			close(clients[user_exit.id].chanel_in)
			delete(clients, user_exit.id)
			user_exit.connection.Close()
		}
	}
}

func handler(connection net.Conn) {
	clientCount++
	user := client{connection, make(chan string), clientCount}

	go clientWriter(user)
	go clientReader(user)
	user.chanel_in <- "-> I connected to: " + connection.LocalAddr().String()
	chanelEntering <- user

}

func main() {
	//lister - серверный сокет
	//listener, err := net.Listen("tcp4", "192.168.0.105:1027")
	listener, err := net.Listen("tcp4", "172.20.10.2:1027")
	if err != nil {
		log.Fatal(err)
		fmt.Println("error", err)
	}
	fmt.Println("<- server created")
	go broadcaster()
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибочка вышла c подключением")
			log.Fatal(err)
		}
		go handler(connection)
	}
}
