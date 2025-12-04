package main

import (
	"log"
)

func main() {
	Config := newConfig()
	application := NewApplication(&Config)
	if err := application.run(); err != nil {
		log.Fatal(err)
	}
}
