package service

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func InitHttpService(configuration *Configuration) {
	router := mux.NewRouter()
	err := http.ListenAndServe(configuration.Address+":"+configuration.Port, router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
