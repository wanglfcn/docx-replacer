package main

import (
	"log"
	"github.com/andlabs/ui"
)

func main() {

	err := ui.Main(mainwindow)
	if err != nil {
		log.Fatalf("fatal error %v", err)
	}
}
