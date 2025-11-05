package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	var (
		id      = r.URL.Query().Get("id")
		service = r.URL.Query().Get("service")
	)

	u.SaveUser(id, service)

	w.Write([]byte("table save success"))
}

func (u *Manage) GetTableHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ipTable, exist := u.ipTable[id]
	if !exist {
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

	r.HandleFunc("GET /", roothHandler)
	r.HandleFunc("GET /health", healthHandler)
	r.HandleFunc("GET /details", detailshHandler)

	r.HandleFunc("GET /table", m.SaveTableHandler)
	r.HandleFunc("GET /table/{id}", m.GetTableHandler)

	srv := &http.Server{
		Addr:         ":80",
		Handler:      r,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 15,
	}

	log.Println("Server has started")

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}

			log.Fatalf("server listening error: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig

	log.Println("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown server error: %v", err)
	}

	log.Println("server exiting")
}
