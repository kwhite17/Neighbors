package neighbors

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/kwhite17/Neighbors/database"
)

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handleCreateNeighbor(w, r)
	case "GET":
		handleGetNeighbor(w, r)
	case "DELETE":
		handleDeleteNeighbor(w, r)
	case "PUT":
	case "PATCH":
		handleUpdateNeighbor(w, r)
	default:
		w.Write([]byte("Invalid Request"))
	}
}

func handleCreateNeighbor(w http.ResponseWriter, r *http.Request) {
	userData := make(map[string]string)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userData)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error decoding user data: %v\n", err)))
		return
	}
	w.Write([]byte("Creating user data for " + userData["username"] + "\n"))
	database.ExecuteWriteQuery(r.Context(), "INSERT into neighbors values (%s)", []interface{}{userData["username"]})
}

func handleGetNeighbor(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimPrefix(r.URL.Path, "/neighbors/")
	log.Println("Fetching user: " + username)
	w.Write([]byte("Hello, neighbor " + username + "\n"))
	result := database.ExecuteReadQuery(r.Context(), "SELECT * from neighbors where username=%s", []interface{}{username})
	if result == nil {
		w.Write([]byte{})
	} else {
		jsonResult, _ := json.Marshal(result)
		w.Write(jsonResult)
	}
}

func handleUpdateNeighbor(w http.ResponseWriter, r *http.Request) {
	userData := make(map[string]string)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userData)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error decoding user data: %v\n", err)))
		return
	}
	w.Write([]byte("Updating user data for " + userData["newName"] + "\n"))
	database.ExecuteWriteQuery(r.Context(), "UPDATE neighbors SET username=%s WHERE username=%s", []interface{}{userData["newName"], userData["oldName"]})
}

func handleDeleteNeighbor(w http.ResponseWriter, r *http.Request) {
	userData := make(map[string]string)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userData)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error decoding user data: %v\n", err)))
		return
	}
	w.Write([]byte("Deleting user data for " + userData["username"] + "\n"))
	database.ExecuteWriteQuery(r.Context(), "DELETE FROM neighbors WHERE username=%s", []interface{}{userData["username"]})

}
