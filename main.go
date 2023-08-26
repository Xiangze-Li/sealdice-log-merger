package main

import (
	"log"
	"os"

	"xiangzeli/logmerger/internal/app"
)

func main() {
	if err := app.App.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
