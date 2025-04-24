package main

import (
	"github.com/joho/godotenv"

	"jcourse_go/internal/infra"
)

func main() {
	_ = godotenv.Load()
	infra.InitDBClient()
	db := infra.GetDBClient()
	err := infra.Migrate(db)
	if err != nil {
		panic(err)
	}
}
