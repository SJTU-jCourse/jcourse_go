package main

import (
	"github.com/joho/godotenv"

	"jcourse_go/config"
	"jcourse_go/internal/app"
)

func main() {
	_ = godotenv.Load()
	c, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	app.Run(c)
}
