package main

import (
	"log"

	flag "github.com/spf13/pflag"

	"jcourse_go/internal/config"
	"jcourse_go/internal/infrastructure/dal"
	"jcourse_go/internal/infrastructure/entity"
)

func main() {
	configPath := flag.StringP("config", "c", "config/config.yaml", "config file path")
	flag.Parse()

	c := config.InitConfig(*configPath)

	db, err := dal.NewPostgresSQL(c.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	err = entity.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}
}
