package main

import (
	"time"
	"fmt"
	"github.com/mmmpa/heaper"
)

func main() {
	go heaper.Run(1, 10)
	for {
		time.Sleep(5 * time.Second)
		rs := heaper.Read()
		for _, r := range rs {
			fmt.Printf("%+v\n", r)
		}
	}
}
