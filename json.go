package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	type errorResult struct {
		Error string `json:"error"`
	}

	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Println("Responding with 5XX error: ", msg)
	}

	respondWithJson(w, code, errorResult{
		Error: msg,
	})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(&payload)
	if err != nil {
		log.Printf("Error marshalling Json: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
	w.Write([]byte("\n"))
}
