package service

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"path"
	"text_search_service/text_search"
)

type Connectors struct {
	textSearch *text_search.TextSearch
}

func InitService(configuration *Configuration) {
	router := mux.NewRouter()
	connectors := initConnectors(configuration)

	router.HandleFunc("/get/{uuid}", get(connectors))

	err := http.ListenAndServe(configuration.Address+":"+configuration.Port, router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func get(connectors *Connectors) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		uuid := path.Base(r.URL.Path)
		log.Println(uuid)
	}
	return http.HandlerFunc(fn)
}

func Close() {
	log.Println("close the service")
}

func initConnectors(configuration *Configuration) *Connectors {
	connectors := &Connectors{}
	return connectors
}
