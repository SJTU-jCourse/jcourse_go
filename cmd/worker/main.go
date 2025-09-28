package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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

	// 2. Listen for signals to gracefully shut down
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	log.Printf("[main] Caught signal: %v. Starting graceful shutdown...", sig)
	err := task.GlobalTaskManager.Shutdown()
	if err != nil {
		log.Printf("[main] Error while shutting down TaskManager: %v\n", err)
	}
	log.Println("[main] Graceful shutdown complete. Exiting.")
	os.Exit(0)

}
