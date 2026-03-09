package main

import (
	"log"

	"github.com/RedPanda13/task_manager/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
