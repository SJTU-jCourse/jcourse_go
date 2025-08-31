package main

import (
	flag "github.com/spf13/pflag"

	"jcourse_go/internal/app"
	"jcourse_go/internal/config"
	"jcourse_go/internal/interface/router"
	"jcourse_go/internal/service"
	"jcourse_go/pkg/util"
)

func main() {
	configPath := flag.StringP("config", "c", "config/config.yaml", "config file path")
	flag.Parse()

	c := config.InitConfig(*configPath)

	if err := util.InitSegWord(); err != nil {
		panic(err)
	}

	err := service.InitLLM()
	if err != nil {
		panic(err)
	}

	s, err := app.NewServiceContainer(c)
	if err != nil {
		panic(err)
	}

	// 3. Start serving
	r := router.RegisterRouter(s)
	_ = r.Run()
}
