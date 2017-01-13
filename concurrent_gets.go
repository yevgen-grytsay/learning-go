package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	urls := []string{
		"http://httpbin.orgs",
		"http://httpbin.org",
		"http://httpbin.org/ip",
	}

	ch := make(chan string)
	for _, url := range urls {
		go MakeRequest(url, ch)
	}
	for range urls {
		fmt.Println(<-ch)
	}
}

func MakeRequest(url string, ch chan<- string) {
	defer func() {
		if err := recover(); err != nil {
			ch <- fmt.Sprintf("Could not fetch url")
		}
	}()
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	ch <- fmt.Sprintf("Body from url '%s' is %d bytes long", url, len(body))
	//if err == nil {
	//    ch <- fmt.Sprintf("Could not fetch url '%s'", url)
	//} else {
	//
	//}
}
