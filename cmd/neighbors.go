package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/kwhite17/Neighbors/pkg/database"
	"github.com/kwhite17/Neighbors/pkg/items"
	"github.com/kwhite17/Neighbors/pkg/login"
	"github.com/kwhite17/Neighbors/pkg/retriever"
	"github.com/kwhite17/Neighbors/pkg/shelters"
)

var NeighborsDatabase database.NeighborsDatasource

func main() {
	dbHost := flag.String("dbhost", "file::memory:?mode=memory&cache=shared", "Name of host on which to run Neighbors")
	developmentMode := flag.Bool("developmentMode", false, "run app in development mode")
	flag.Parse()
	log.Println("Connecting to host", *dbHost)
	log.Println("Development mode set to:", *developmentMode)

	NeighborsDatabase = database.NeighborsDatasource{Database: database.InitDatabase(*dbHost, *developmentMode)}
	shelterManager := &shelters.ShelterManager{Datasource: NeighborsDatabase}
	itemManager := &items.ItemManager{Datasource: NeighborsDatabase}
	shelterSessionManager := &login.ShelterSessionManager{Datasource: NeighborsDatabase}
	mux := http.NewServeMux()

	mux.Handle("/shelters/", buildShelterServiceHandler(shelterManager, shelterSessionManager, itemManager))
	mux.Handle("/items/", buildItemServiceHandler())
	mux.Handle("/session/", buildLoginServiceHandler(shelterManager, shelterSessionManager))
	mux.HandleFunc("/", loadHomePage)
	http.ListenAndServe(":8080", mux)
}

func loadHomePage(w http.ResponseWriter, r *http.Request) {
	t, err := retriever.RetrieveTemplate("home/index")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	t.Execute(w, nil)
}

func buildShelterServiceHandler(shelterManager *shelters.ShelterManager, shelterSessionManager *login.ShelterSessionManager, itemManager *items.ItemManager) shelters.ShelterServiceHandler {
	return shelters.ShelterServiceHandler{
		ShelterRetriever:      &shelters.ShelterRetriever{},
		ShelterManager:        shelterManager,
		ShelterSessionManager: shelterSessionManager,
		ItemManager:           itemManager,
	}
}

func buildItemServiceHandler() items.ItemServiceHandler {
	return items.ItemServiceHandler{
		ItemRetriever: &items.ItemRetriever{},
		ItemManager:   &items.ItemManager{Datasource: NeighborsDatabase},
	}
}

func buildLoginServiceHandler(shelterManager *shelters.ShelterManager, shelterSessionManager *login.ShelterSessionManager) login.LoginServiceHandler {
	return login.LoginServiceHandler{
		ShelterSessionManager: shelterSessionManager,
		LoginRetriever:        &login.LoginRetriever{},
	}
}
