package main

import (
	"fmt"
	"log"
	"net"
)

type client struct {
	connection net.Conn
	chanel     chan string
}

var (
	entering = make(chan string)
	leaving  = make(chan string)
	message  = make(chan string)
)

func clientWriter(connection net.Conn, chanel chan string) {
	fmt.Println("\nfunc client writer")
}

// широковещатель
func broadcaster() {
	clients := make([]client, 0)
	for {
		select {
		case msg := <-entering:
			fmt.Println("\n case entering")
			clients = append(clients)
			clients[len(clients)-1].msg_input <- "\nhi #" + clients[len(clients)-1].conn.RemoteAddr().String()
			fmt.Println(<-clients[len(clients)-1].msg_input)

			for i := range clients {
				clients[i].msg_input <- "\nnew user"
			}
		default:
			fmt.Println("\n case default")
			for i := range clients {
				clients[i].msg_input <- "\nnothing"
			}
		}
	}
}

func handler(connection net.Conn) {
	defer connection.Close()

	chanel := make(chan string)

	go clientWriter(connection, chanel)

	chanel <- "I'am " + connection.LocalAddr().String()
	entering <- client{conneconnection, chanel}

	msg := make([]byte, 80)
	for {
		_, err := connection.Read(msg)
		if err != nil {
			break
		}
		fmt.Println(string(msg))
	}

}

func main() {
	//lister - серверный сокет
	listener, err := net.Listen("tcp4", "192.168.0.105:1027")
	if err != nil {
		log.Fatal(err)
		fmt.Println("error", err)
	}
	fmt.Println(listener)
	fmt.Println("ok")
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
