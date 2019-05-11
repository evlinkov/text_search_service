package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"text_search_service/text_search"
	"text_search_service/util"
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

type WordResponse struct {
	Text       string    `json:"text"`
	Uuid       uuid.UUID `json:"uuid"`
	Popularity int64     `json:"popularity"`
}

func (service *Service) InitService(configuration *Configuration) {
	router := mux.NewRouter()
	service.configuration = configuration
	service.initConnectors()

	router.HandleFunc("/get/{uuid}", get(service))

	go func() {
		err := http.ListenAndServe(service.configuration.Address+":"+service.configuration.Port, router)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
}

func get(service *Service) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		uuid, err := util.ParseStringToUUID(path.Base(r.URL.Path))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("bad request %v\n", err)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		word, ok := service.connectors.textSearch.GetWordByUUID(uuid)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(convert(word))
		}
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

func convert(word text_search.Word) WordResponse {
	wordResponse := WordResponse{}
	wordResponse.Text = word.Text
	wordResponse.Uuid = word.Uuid
	wordResponse.Popularity = word.Popularity
	return wordResponse
}
