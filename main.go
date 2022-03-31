package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

var ClientCount int = 0

func Handler(connection net.Conn) {
	defer connection.Close()
	fmt.Println("\nConnection #", connection)
	ClientCount++
	connection.Write([]byte("\nHi! You are #" + strconv.Itoa(ClientCount)))
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

	connection := make([]net.Conn, 0)
	for {
		client, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибочка вышла")
			log.Fatal(err)
		}
		connection = append(connection, client)
		fmt.Println("\namount connection ", len(connection))
		go Handler(client)
		for i := 0; i < len(connection)-1; i++ {
			fmt.Println("\n i =", i)
			connection[i].Write([]byte("\nnew client #" + string(ClientCount)))
		}
	}

}
