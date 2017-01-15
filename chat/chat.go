package main

import (
	"net"
	"container/list"
	"bufio"
	"fmt"
)

var members = list.New()
var chat = make(chan message)

func main() {
	ln, err := net.Listen("tcp", ":30333")
	if err != nil {
		panic(err)
	}
	connGenerator := connectionGenerator(ln)
	nickGenerator := nickGenerator()
	go run()
	for con := range connGenerator {
		go handleConnection(con, <-nickGenerator)
	}
}

type message struct {
	Nick	 string
	UserChan chan string
	Text     string
}

func run() {
	for {
		select {
		case msg := <-chat:
			text := fmt.Sprintf("%s: %s", msg.Nick, msg.Text)
			for e := members.Front(); e != nil; e = e.Next() {
				if e.Value != msg.UserChan {
					e.Value.(chan string) <- text
				}
			}
		}
	}
}

func handleConnection(conn net.Conn, nick string) {
	in := make(chan string)
	members.PushBack(in)
	userInput := userInputGenerator(conn)
	for {
		select {
		case msg := <-in:
			conn.Write([]byte(msg))
		case msg := <-userInput:
			chat <- message{nick, in, msg}
		}
	}
}

func userInputGenerator(conn net.Conn) chan string {
	b := bufio.NewReader(conn)
	ch := make(chan string)
	go func() {
		for {
			if line, err := b.ReadString('\n'); err == nil {
				ch <- line
			}
		}
	}()
	return ch
}

func connectionGenerator(ln net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	go func () {
		for {
			if conn, err := ln.Accept(); err == nil {
				ch <- conn
			}
		}
	}()
	return ch
}

func nickGenerator() chan string {
	ch := make(chan string)
	go func() {
		i := 1
		for {
			ch <- fmt.Sprintf("User %d", i)
			i += 1
		}
	}()
	return ch
}
