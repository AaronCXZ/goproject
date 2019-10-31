package main

import (
	"fmt"
	"log"

	"github.com/robfig/cron"
)

func main() {
	spec := "*/5 * * * * *"
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc(spec, callFunc)
	if err != nil {
		fmt.Println(err)
	}
	c.Start()
	select {}
}

func callFunc() {
	log.Println("Hello World!")
}
