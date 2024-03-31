package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	website := []string{ // Arg this
		"http://amazon.com",
		"http://google.com",
	}

	c := make(chan string)

	for _, w := range website {
		go checkLink(w, c) // Create a separate goroutine
	}
	
	for w := range c { // Loop over channels, infinite loop
		go func(website string) { 
			time.Sleep(time.Second * 5)
			checkLink(website, c)
		}(w) // Function literal, copied in memory so that its stable for the function
	}

}

func checkLink(website string, c chan string) { // Declare the channel here
	res, err := http.Get(website)
	if err != nil {
		fmt.Println(website, "is down!")
		c <- website
		return
	}
	fmt.Println(website, res.StatusCode) // JSON for Prometheus?
	c <- website
}
