package main

import (
	"net"
	"container/list"
	"bufio"
)

var members = list.New()
var chat = make(chan []byte)

func main() {
	ln, err := net.Listen("tcp", ":30333")
	if err != nil {
		panic(err)
	}
	connGenerator := connectionGenerator(ln)
	go run()
	for con := range connGenerator {
		go handleConnection(con)
	}
}

func run() {
	for {
		select {
		case msg := <-chat:
			for e := members.Front(); e != nil; e = e.Next() {
				if user, ok := e.Value.(chan []byte); ok {
					user <- msg
				}
			}
		}
	}
}

func handleConnection(conn net.Conn) {
	in := make(chan []byte)
	members.PushBack(in)
	userInput := userInputGenerator(conn)
	for {
		select {
		case msg := <-in:
			conn.Write(msg)
		case msg := <-userInput:
			chat <- msg
		}
	}
}

func userInputGenerator(conn net.Conn) chan []byte{
	b := bufio.NewReader(conn)
	ch := make(chan []byte)
	go func() {
		for {
			if line, err := b.ReadBytes('\n'); err == nil {
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
