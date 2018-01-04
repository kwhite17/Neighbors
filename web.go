package main

import "net/http"
import "github.com/kwhite17/Neighbors/neighbors"
import "github.com/kwhite17/Neighbors/samaritans"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/neighbors/", neighbors.RequestHandler)
	mux.HandleFunc("/samaritans/", samaritans.RequestHandler)
	http.ListenAndServe(":8080", mux)
}
