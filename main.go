package main

import (
	"fmt"
	"log"
	"songlib/confreader"
)

func main() {
	config, err := confreader.LoadConfig()
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("%+v", config)
}
