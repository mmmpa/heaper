package main

import (
	"mmmpa.net/heaper"
	"time"
	"fmt"
)

func main() {
	for {
		go heaper.Run(1, 10)
		time.Sleep(5 * time.Second)
		rs := heaper.Read()

		for _, r := range rs {
			fmt.Printf("%+v\n", r)
		}
	}
}
