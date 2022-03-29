package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func connection(client net.Listener) {
	// ожидаю новые подключения "слушать порт"
	connection, err := client.Accept()
	if err != nil {
		log.Fatal(err)
	}
	timeout := time.Now().Add(time.Second)
	connection.SetReadDeadline(timeout)

	fmt.Println("\nConnection #", connection)

	// defer - выполнить по завершении программы
	// направить пакет клиенту о заверешении соединения
	defer connection.Close()
}

func sendMsg(client net.Conn, msg string) {
	msg_out := []byte("Prvet Misha")
	client.Write(msg_out)
}

func getMsg(client net.Conn) string {
	msg := make([]byte, 80)
	_, err := client.Read(msg)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(msg))
	return string(msg)
}

func main() {
	//lister - серверный сокет
	listener, err := net.Listen("tcp4", "192.168.1.62:1027")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(listener)
}
