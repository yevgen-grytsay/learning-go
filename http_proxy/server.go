package main

import (
	"net"
	"net/http"
	"bufio"
	"fmt"
	"net/url"
)

var client = http.Client{}

func main() {
	ln, err := net.Listen("tcp", ":30333")
	if err != nil {
		panic(err)
	}
	connGenerator := connectionGenerator(ln)
	for con := range connGenerator {
		go handleConnection(con)
	}
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

func handleConnection(conn net.Conn) {
	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		panic(err)
	}
	fmt.Println(req)
	u, err := url.Parse(req.RequestURI)
	if err != nil {
		panic(err)
	}
	req.URL = u
	req.RequestURI = ""

	response, err := client.Do(req)
	if err == nil {
		response.Write(conn)
	} else {
		fmt.Println(err)
		conn.Write([]byte("HTTP/1.1 500 Cannot reach destination\r\n\r\n"))
	}
	conn.Close()
}