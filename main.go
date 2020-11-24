package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	err := os.Mkdir("streams", os.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatalf("%+v", err)
	}
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	// public site
	fpub := http.FileServer(http.Dir("./public"))
	r.PathPrefix("/players/").Handler(http.StripPrefix("/players/", fpub))

	// stream files
	fstream := http.FileServer(http.Dir("./streams"))
	r.PathPrefix("/watch/").Handler(http.StripPrefix("/watch/", fstream))

	// stream endpoint
	r.HandleFunc("/publish/{streamid}", PublishHandle)
	log.Print("HTTP server listening on 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}

// PublishHandle handles accepting the stream from ffmpeg and storing it locally to disk
func PublishHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		os.Remove("streams/" + mux.Vars(r)["streamid"])
	}
	if r.Method != http.MethodPut {
		return
	}
	// Only accept PUT for adding new stream files
	defer r.Body.Close()
	file, err := os.Create("streams/" + mux.Vars(r)["streamid"])
	defer file.Close()
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(file, r.Body)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] - %s", r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
