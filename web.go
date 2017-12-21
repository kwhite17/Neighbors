package main

import "net/http"
import "github.com/kwhite17/Neighbors/neighbors"

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World\n"))
}

func main() {
	http.HandleFunc("/neighbors/", neighbors.RequestHandler)
	http.ListenAndServe(":8080", nil)
}
