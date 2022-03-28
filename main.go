package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	// local message
	//lister - серверный сокет
	listener, err := net.Listen("tcp4", "192.168.1.62:1027")
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
}
