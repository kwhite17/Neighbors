package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/kwhite17/Neighbors/pkg/managers"
	"github.com/kwhite17/Neighbors/pkg/resources"
	"github.com/kwhite17/Neighbors/pkg/retrievers"

	"github.com/kwhite17/Neighbors/pkg/database"
	"github.com/kwhite17/Neighbors/pkg/email"
	"github.com/sendgrid/sendgrid-go"
	"gopkg.in/gomail.v2"
)

type EnvironmentConfig struct {
	Datasource  database.Datasource
	EmailSender email.EmailSender
}

func BuildDatasource(driver string, host string, developmentMode bool) database.Datasource {
	return database.BuildDatasource(driver, host, developmentMode)
}

func BuildEmailSender(developmentMode bool, userManager *managers.UserManager) email.EmailSender {
	if developmentMode {
		return &email.LocalSender{Dialer: &gomail.Dialer{Host: "localhost", Port: 25}, UserManager: userManager}
	}
	return &email.SendGridSender{Client: sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY")), UserManager: userManager}
}

func main() {
	driver := flag.String("dbDriver", "sqlite3", "Name of database driver to use")
	dbHost, dbHostFound := os.LookupEnv("DATABASE_URL")
	developmentMode := flag.Bool("developmentMode", false, "run app in development mode")
	if !dbHostFound {
		dbHost = "file::memory:?mode=memory&cache=shared"
	}
	port, portFound := os.LookupEnv("PORT")
	if !portFound {
		port = "8080"
	}
	flag.Parse()
	log.Println("Connecting to host", dbHost)
	log.Println("Development mode set to:", *developmentMode)

	environment := &EnvironmentConfig{}
	environment.Datasource = BuildDatasource(*driver, dbHost, *developmentMode)
	userManager := &managers.UserManager{Datasource: environment.Datasource}
	itemManager := &managers.ItemManager{Datasource: environment.Datasource}
	userSessionManager := &managers.UserSessionManager{Datasource: environment.Datasource}
	environment.EmailSender = BuildEmailSender(*developmentMode, userManager)

	mux := http.NewServeMux()
	mux.Handle("/shelters/", buildUserServiceHandler(userManager, userSessionManager, itemManager))
	mux.Handle("/items/", buildItemServiceHandler(userSessionManager, itemManager, userManager, environment))
	mux.Handle("/session/", buildLoginServiceHandler(userManager, userSessionManager))
	mux.Handle("/", buildHomeServiceHandler(userSessionManager))
	http.ListenAndServe(":"+port, mux)
}

func buildHomeServiceHandler(userSessionManager *managers.UserSessionManager) resources.HomeServiceHandler {
	return resources.HomeServiceHandler{
		UserSessionManager: userSessionManager,
	}
}

func buildUserServiceHandler(userManager *managers.UserManager, userSessionManager *managers.UserSessionManager, itemManager *managers.ItemManager) resources.UserServiceHandler {
	return resources.UserServiceHandler{
		UserRetriever:      &retrievers.ShelterRetriever{},
		UserManager:        userManager,
		UserSessionManager: userSessionManager,
		ItemManager:        itemManager,
	}
}

func buildItemServiceHandler(userSessionManager *managers.UserSessionManager, itemManager *managers.ItemManager, userManager *managers.UserManager, environment *EnvironmentConfig) resources.ItemServiceHandler {
	return resources.ItemServiceHandler{
		ItemRetriever:      &retrievers.ItemRetriever{},
		ItemManager:        itemManager,
		UserSessionManager: userSessionManager,
		EmailSender:        environment.EmailSender,
	}
}

func buildLoginServiceHandler(userManager *managers.UserManager, userSessionManager *managers.UserSessionManager) resources.LoginServiceHandler {
	return resources.LoginServiceHandler{
		UserManager:        userManager,
		UserSessionManager: userSessionManager,
		LoginRetriever:     &retrievers.LoginRetriever{},
	}
}
