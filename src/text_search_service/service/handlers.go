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
	"sort"
	"strings"
	"text_search_service/text_search"
	"text_search_service/util"
)

const (
	MaxLimitSearchResponse         = 10
	NumberOfSymbolsAroundSubstring = 2
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

type WordResponseSearch struct {
	Preview    string    `json:"preview"`
	Uuid       uuid.UUID `json:"uuid"`
	Popularity int64     `json:"popularity"`
}

type SearchResponse struct {
	Words []WordResponseSearch `json:"results"`
}

type Request struct {
	Query string `json:"query"`
}

type SaveResponse struct {
	Id string `json:"id"`
}

func (service *Service) InitService(configuration *Configuration) {
	router := mux.NewRouter()
	service.configuration = configuration
	service.initConnectors()

	router.HandleFunc("/get/{uuid}", get(service))
	router.HandleFunc("/save", save(service))
	router.HandleFunc("/search", search(service))

	go func() {
		err := http.ListenAndServe(service.configuration.Address+":"+service.configuration.Port, router)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
}

func save(service *Service) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("error while reading body %v", err), http.StatusInternalServerError)
			log.Printf("error while reading body %v\n", err)
			return
		}

		request := Request{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response := SaveResponse{}
		response.Id = fmt.Sprintf("%v", service.connectors.textSearch.AddWord(request.Query))
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
	return http.HandlerFunc(fn)
}

func search(service *Service) http.HandlerFunc {
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

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("error while reading body %v", err), http.StatusInternalServerError)
			log.Printf("error while reading body %v\n", err)
			return
		}

		request := Request{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		words := service.connectors.textSearch.Search(request.Query)
		sort.Sort(ByPopularity(words))
		if len(words) > MaxLimitSearchResponse {
			words = words[:MaxLimitSearchResponse]
		}
		response := SearchResponse{}
		response.Words = make([]WordResponseSearch, 0, len(words))
		for _, word := range words {
			wordSearchResponse := WordResponseSearch{}
			wordSearchResponse.Popularity = word.Popularity
			wordSearchResponse.Uuid = word.Uuid
			index := strings.Index(word.Text, request.Query)
			if index == -1 {
				log.Printf("ERROR in tree, request : %s\n", request.Query)
				continue
			}
			isBefore := ""
			before := ""
			if index > NumberOfSymbolsAroundSubstring {
				isBefore = "..."
				before = word.Text[(index - NumberOfSymbolsAroundSubstring):index]
			} else {
				before = word.Text[0:index]
			}

			isAfter := ""
			after := ""
			if len(word.Text)-(index+len(request.Query)) > NumberOfSymbolsAroundSubstring {
				isAfter = "..."
				from := index + len(request.Query)
				after = word.Text[from : from+NumberOfSymbolsAroundSubstring]
			} else {
				after = word.Text[(index + len(request.Query)):]
			}
			wordSearchResponse.Preview = fmt.Sprintf("%s%s%s%s%s", isBefore, before, request.Query, after, isAfter)
			response.Words = append(response.Words, wordSearchResponse)
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
	return http.HandlerFunc(fn)
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

type ByPopularity []text_search.Word

func (a ByPopularity) Len() int      { return len(a) }
func (a ByPopularity) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPopularity) Less(i, j int) bool {
	return a[i].Popularity > a[j].Popularity
}
