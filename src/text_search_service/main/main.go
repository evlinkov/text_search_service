package main

import (
	"text_search_service/service"
)

func main() {
	configuration := service.GetConfiguration()
	service.InitHttpService(configuration)
}
