package main

import (
	"log"
	"net/http"

	"audioinsights-service-go/internal/httpapi"
	"audioinsights-service-go/internal/spotify"
)

const serverAddress = ":8080"

func main() {
	server := httpapi.NewServer(spotify.NewClient())

	log.Printf("audio insights service listening on %s", serverAddress)
	if err := http.ListenAndServe(serverAddress, server); err != nil {
		log.Fatal(err)
	}
}
