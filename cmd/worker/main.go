package main

import (
	flag "github.com/spf13/pflag"

	"jcourse_go/internal/app"
	"jcourse_go/internal/config"
	"jcourse_go/internal/interface/task"
)

func main() {
	configPath := flag.StringP("config", "c", "config/config.yaml", "config file path")
	flag.Parse()

	c := config.InitConfig(*configPath)

	s, err := app.NewServiceContainer(c)
	if err != nil {
		panic(err)
	}

	srv, mux := task.NewAsyncTaskServer(s)
	if err = srv.Run(mux); err != nil {
		panic(err)
	}
}
