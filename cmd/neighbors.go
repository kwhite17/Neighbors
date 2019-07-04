package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/kwhite17/Neighbors/pkg/database"
	"github.com/kwhite17/Neighbors/pkg/managers"
	"github.com/kwhite17/Neighbors/pkg/resources"
	"github.com/kwhite17/Neighbors/pkg/retrievers"
)

var NeighborsDatabase database.NeighborsDatasource

func main() {
	dbHost := flag.String("dbhost", "file::memory:?mode=memory&cache=shared", "Name of host on which to run Neighbors")
	developmentMode := flag.Bool("developmentMode", false, "run app in development mode")
	flag.Parse()
	log.Println("Connecting to host", *dbHost)
	log.Println("Development mode set to:", *developmentMode)

	NeighborsDatabase = database.NeighborsDatasource{Database: database.InitDatabase(*dbHost, *developmentMode)}
	shelterManager := &managers.ShelterManager{Datasource: NeighborsDatabase}
	itemManager := &managers.ItemManager{Datasource: NeighborsDatabase}
	shelterSessionManager := &managers.ShelterSessionManager{Datasource: NeighborsDatabase}
	mux := http.NewServeMux()

	mux.Handle("/shelters/", buildShelterServiceHandler(shelterManager, shelterSessionManager, itemManager))
	mux.Handle("/items/", buildItemServiceHandler())
	mux.Handle("/session/", buildLoginServiceHandler(shelterManager, shelterSessionManager))
	mux.HandleFunc("/", loadHomePage)
	http.ListenAndServe(":8080", mux)
}

func loadHomePage(w http.ResponseWriter, r *http.Request) {
	t, err := retrievers.RetrieveTemplate("home/index")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	t.Execute(w, nil)
}

func buildShelterServiceHandler(shelterManager *managers.ShelterManager, shelterSessionManager *managers.ShelterSessionManager, itemManager *managers.ItemManager) resources.ShelterServiceHandler {
	return resources.ShelterServiceHandler{
		ShelterRetriever:      &retrievers.ShelterRetriever{},
		ShelterManager:        shelterManager,
		ShelterSessionManager: shelterSessionManager,
		ItemManager:           itemManager,
	}
}

func buildItemServiceHandler() resources.ItemServiceHandler {
	return resources.ItemServiceHandler{
		ItemRetriever: &retrievers.ItemRetriever{},
		ItemManager:   &managers.ItemManager{Datasource: NeighborsDatabase},
	}
}

func buildLoginServiceHandler(shelterManager *managers.ShelterManager, shelterSessionManager *managers.ShelterSessionManager) resources.LoginServiceHandler {
	return resources.LoginServiceHandler{
		ShelterManager:        shelterManager,
		ShelterSessionManager: shelterSessionManager,
		LoginRetriever:        &retrievers.LoginRetriever{},
	}
}
