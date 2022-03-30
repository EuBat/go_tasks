package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// сделать переменные сеттеров cosnt
// сделать геттеры const функцией
// где деструктор :(
// конструктора нет. деструктура видимо тоже нет.
var wg sync.WaitGroup

type Client struct {
	Connection net.Conn
	status     bool
}

func (c *Client) SendMsg(message string) {
	msg := []byte(message)
	c.Connection.Write(msg)
}

func (c *Client) GetMsg() string {

	msg := make([]byte, 80)

	go c.Connection.Read(msg)

	if string(msg) == "bye" {
		c.Connection.Close()
	}
	fmt.Println(string(msg))
	return string(msg)
}

func (c *Client) SetTimeout(second time.Duration) {
	timeout := time.Now().Add(second + time.Second)
	c.Connection.SetReadDeadline(timeout)
}
func (c *Client) Connect(server net.Listener) {
	var err error
	c.Connection, err = server.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Connection.Close()
	//c.SetTimeout(time.Second)
	c.SendMsg("Hi Eugen")
	c.GetMsg()
	fmt.Println("\nConnection #", c.Connection)
}

func (c *Client) ConnectClose() {
	c.Connection.Close()
	c.SendMsg("bye-bye")
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

	clients := make([]Client, 10)

	wg.Add(5)
	go clients[0].Connect(listener)
	go clients[1].Connect(listener)
	go clients[2].Connect(listener)
	go clients[3].Connect(listener)
	go clients[4].Connect(listener)

	// for i := range clients {
	// 	fmt.Println("/n#G:", i)
	// 	fmt.Println("/n#:", clients[i].Connection)
	// 	defer wg.Done()
	// }

	// for i := range clients {
	// 	go clients[i].SendMsg("bye-bye")
	// 	go clients[i].ConnectClose()
	// }

	wg.Wait()
}
