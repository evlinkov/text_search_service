package main

import (
	"os"
	"os/signal"
	"syscall"
	"text_search_service/service"
)

var (
	interrupt chan os.Signal
)

func main() {
	configuration := service.GetConfiguration()
	interrupt = make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	service := &service.Service{}
	service.InitService(configuration)
	select {
	case <-interrupt:
		service.Close()
		break
	}
}
