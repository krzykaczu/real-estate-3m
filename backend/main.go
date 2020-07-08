package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type AdData struct {
	Id    string  `json:"id"`
	Url   string  `json:"url"`
	Loc   string  `json:"loc"`
	Price string  `json:"price"`
	Rooms string  `json:"rooms"`
	M2    string  `json:"m2"`
	Floor string  `json:"floor"`
}

type AdDataNodes struct {
	Data []AdData `json:"data"`
}

type SingleAttr struct {
	ObjectID   string
	SingleLine string
}

type Attributes struct {
	Attributes SingleAttr
}

type Records struct {
	Records []Attributes
}

func get(w http.ResponseWriter, r *http.Request) {
	filename := "./assets/scrap3m.json"
	file, _ := ioutil.ReadFile(filename)
	data := AdDataNodes{}
	_ = json.Unmarshal([]byte(file), &data)

	token := "EfpuXlCvYMNSZ1Vbn5EZB9jfbNbqvvDqToqfqnq7w6Xlys3gKZV8aE8ixOU2S1NEFsA7aXqZEJO8uaGD66s6RsWeY6WX75WJaobU5WUOfgpc-qFDo-7ntHiX7b03UejUnAiw8ZDjLUNn9H2kcZtSZ73bwVAxtIv1VQc52Nz3CDEMHuJiFGV54pt8aElrxzGXfqw3LAKG0RpAQCBIRBmZ96G5kYr52VA5Dv55sqEcu98."
	var records Records
	for i, v := range data.Data {
		singleAttr := SingleAttr{
			ObjectID: strconv.Itoa(i + 1),
			SingleLine: v.Loc,
		}
		attr := Attributes{singleAttr}
		records.Records = append(records.Records, attr)
	}
	marshaledRecords, err := json.Marshal(records)
	if err != nil {
			fmt.Println(err)
			return
		}
 	encodedRecords := url.QueryEscape(string(marshaledRecords))
	URL := "https://geocode.arcgis.com/arcgis/rest/services/World/GeocodeServer/geocodeAddresses?addresses=" + encodedRecords + "&f=pjson&token=" + token

	fmt.Println(URL)

	//resp, err := http.Get("http://example.com/")

	b, err := json.Marshal(data.Data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "error marshalling data"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "not found"}`))
}

func main() {
	//file, err := os.Open(filename)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//defer file.Close()


	//file, _ := ioutil.ReadFile(filename)
	//data := AdDataNodes{}
	//_ = json.Unmarshal([]byte(file), &data)
	//fmt.Println(data)

	r := mux.NewRouter()
	r.HandleFunc("/data", get).Methods(http.MethodGet)
	r.HandleFunc("", notFound)
	log.Fatal(http.ListenAndServe(":8080", r))
}