package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Iowel/kube-dep-srv/details"
)

type Manage struct {
	ipTable map[string]string
}

func NewManage() *Manage {
	return &Manage{
		ipTable: make(map[string]string),
	}
}

func (u *Manage) SaveUser(id, service string) {
	u.ipTable[id] = service
}

func (u *Manage) SaveTableHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	service := r.URL.Query().Get("service")

	u.SaveUser(id, service)

	w.Write([]byte("table save success"))
}

func (u *Manage) GetTableHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ipTable, ok := u.ipTable[id]
	if !ok {
		http.Error(w, "table not exist", http.StatusBadRequest)
	}

	w.Write([]byte(ipTable))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("checking application health")

	response := map[string]string{
		"status":    "UP",
		"timestamp": time.Now().String(),
	}

	json.NewEncoder(w).Encode(response)
}

func roothHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving the homepage")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Application is up running")
}

func detailshHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching the details")

	hostname, err := details.GetHostName()
	if err != nil {
		panic(err)
	}

	ip, err := details.GetIP()
	if err != nil {
		panic(err)
	}

	response := map[string]string{
		"hostname": hostname,
		"ip":       ip.String(),
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	r := http.NewServeMux()

	m := NewManage()

	r.HandleFunc("/", roothHandler)
	r.HandleFunc("/health", healthHandler)
	r.HandleFunc("/details", detailshHandler)

	r.HandleFunc("/table", m.SaveTableHandler)
	r.HandleFunc("/table/{id}", m.GetTableHandler)

	s := &http.Server{
		Addr:    ":80",
		Handler: r,
	}

	log.Println("Server has started")

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
