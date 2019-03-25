package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kwhite17/Neighbors/pkg/database"
	"github.com/kwhite17/Neighbors/pkg/items"
	"github.com/kwhite17/Neighbors/pkg/users"
)

func main() {
	directory, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	directory, err = filepath.EvalSymlinks(directory)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(directory)
	mux := http.NewServeMux()
	mux.Handle("/users/", users.UserServiceHandler{Database: database.NeighborsDatabase})
	mux.Handle("/items/", items.ItemResourceHandler{Database: database.NeighborsDatabase})
	mux.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir(directory+"/templates/"))))

	http.ListenAndServe(":8080", mux)
}
