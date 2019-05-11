package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"text_search_service/text_search"
)

const (
	DELIMITER = "±"
)

type Connectors struct {
	textSearch *text_search.TextSearch
}

type Service struct {
	connectors    *Connectors
	configuration *Configuration
}

func (service *Service) InitService(configuration *Configuration) {
	router := mux.NewRouter()
	service.configuration = configuration
	service.initConnectors()

	router.HandleFunc("/get/{uuid}", get())

	go func() {
		err := http.ListenAndServe(service.configuration.Address+":"+service.configuration.Port, router)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
}

func get() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		uuid := path.Base(r.URL.Path)
		log.Println(uuid)
	}
	return http.HandlerFunc(fn)
}

func (service *Service) Close() {
	log.Println("try close the service")
	service.printData()
}

func (service *Service) initConnectors() {
	connectors := &Connectors{}
	connectors.textSearch = text_search.InitTextSearch(service.readPreviousData())
	service.connectors = connectors
}

func (service *Service) readPreviousData() []text_search.Word {
	data, err := ioutil.ReadFile(service.configuration.Filename)
	if err != nil {
		fmt.Printf("error reading file %v\n", err)
		return nil
	}
	words := make([]text_search.Word, 0)
	err = json.Unmarshal(data, &words)
	if err != nil {
		fmt.Printf("error unmarshal file %v\n", err)
		return nil
	}
	return words
}

func (service *Service) printData() {
	file, err := os.Create(service.configuration.Filename)
	if err != nil {
		log.Printf("error while creating file %v\n", err)
		return
	}
	defer file.Close()
	words := service.connectors.textSearch.GetAllWords()
	data, err := json.Marshal(words)
	if err != nil {
		log.Printf("error marshal data %+v\n", data)
		return
	}
	file.Write(data)
}
