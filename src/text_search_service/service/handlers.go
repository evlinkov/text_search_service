package service

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func InitService(configuration *Configuration) {
	router := mux.NewRouter()
	err := http.ListenAndServe(configuration.Address+":"+configuration.Port, router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func Close() {
	log.Println("close the service")
}
