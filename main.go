package main

import (
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
	clientCount int = 0
	entering        = make(chan client)
	exit            = make(chan client)
	message         = make(chan string)
)

// печать входящих сообщений
func clientWriter(c client) {
	for {
		c.connection.Write([]byte(<-c.chanel_in))
	}
}

// чтение из консоли исходящих сообщений
func clientReader(c client) {
	msg := make([]byte, 4)
	for {
		_, err := c.connection.Read(msg)
		if err != nil {
			break
		}
		fmt.Println(string(msg))
		if string(msg) == "exit" {
			fmt.Println("exitif")
			exit <- c
		} else {
			message <- string(msg)
		}
	}
}

// широковещатель
func broadcaster() {
	clients := make(map[int]client)
	//clients := make([]client, 0)
	for {
		select {
		case user := <-entering:
			//add new user to slice
			clients[user.id] = user

			//meeting with the new user
			clients[clientCount].chanel_in <- "<- hi #" + clients[clientCount].connection.RemoteAddr().String() + "\n"

			// tell everybody about new user
			for i := range clients {
				clients[i].chanel_in <- "\n new user " + user.connection.RemoteAddr().String() + " connected\n"
			}
		case user_msg := <-message:
			msg := strings.Split(user_msg, " ")
			userId, _ := strconv.Atoi(msg[0])
			for i := range clients {
				if clients[i].id == userId {
					clients[i].chanel_in <- msg[1]
				}
			}
		case user := <-exit:
			delete(clients, user.id)
			clientCount--
			close(user.chanel_in)
			user.connection.Close()
			//fmt.Println("\n case default")
		}
	}
}

func handler(connection net.Conn) {
	clientCount++
	user := client{connection, make(chan string), clientCount}

	go clientWriter(user)
	go clientReader(user)
	user.chanel_in <- "-> I connected to: " + connection.LocalAddr().String() + "\n"
	entering <- user

}

func main() {
	//lister - серверный сокет
	//listener, err := net.Listen("tcp4", "192.168.0.105:1027")
	listener, err := net.Listen("tcp4", "172.20.10.2:1027")
	if err != nil {
		log.Fatal(err)
		fmt.Println("error", err)
	}
	fmt.Println("server created")
	go broadcaster()
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибочка вышла c подключением")
			log.Fatal(err)
		}
		fmt.Println(connection.LocalAddr().String())
		go handler(connection)
	}
}
