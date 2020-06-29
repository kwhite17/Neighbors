package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/kwhite17/Neighbors/pkg/managers"
	"github.com/kwhite17/Neighbors/pkg/resources"
	"github.com/kwhite17/Neighbors/pkg/retrievers"
	"gopkg.in/gomail.v2"

	"github.com/kwhite17/Neighbors/pkg/database"
	"github.com/kwhite17/Neighbors/pkg/email"
	"github.com/sendgrid/sendgrid-go"
)

type EnvironmentConfig struct {
	Datasource  database.Datasource
	EmailSender email.EmailSender
}

func buildDatasource(driver string, host string, developmentMode bool) database.Datasource {
	return database.BuildDatasource(driver, host, developmentMode)
}

func buildEmailSender(developmentMode bool, userManager *managers.UserManager) email.EmailSender {
	if developmentMode {
		return &email.LocalSender{Dialer: &gomail.Dialer{Host: "localhost", Port: 25}, UserManager: userManager}
	}
	return &email.SendGridSender{Client: sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY")), UserManager: userManager}
}

func buildEnvironment(datasource database.Datasource, emailSender email.EmailSender) *EnvironmentConfig {
	return &EnvironmentConfig{Datasource: datasource, EmailSender: emailSender}
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

	datasource := buildDatasource(*driver, dbHost, *developmentMode)
	userManager := &managers.UserManager{Datasource: datasource}
	itemManager := &managers.ItemManager{Datasource: datasource}
	userSessionManager := &managers.UserSessionManager{Datasource: datasource}
	environment := buildEnvironment(datasource, buildEmailSender(*developmentMode, userManager))

	router := mux.NewRouter()
	router.PathPrefix("/shelters").Handler(buildUserServiceHandler(userSessionManager, userManager, itemManager))
	router.PathPrefix("/items").Handler(buildItemServiceHandler(userSessionManager, itemManager, environment))
	router.PathPrefix("/session").Handler(buildLoginServiceHandler(userSessionManager, userManager))
	router.PathPrefix("/").Handler(buildHomeServiceHandler(userSessionManager))
	http.ListenAndServe(":"+port, router)
}

func buildHomeServiceHandler(userSessionManager *managers.UserSessionManager) resources.HomeServiceHandler {
	return resources.HomeServiceHandler{
		UserSessionManager: userSessionManager,
	}
}

func buildUserServiceHandler(userSessionManager *managers.UserSessionManager, userManager *managers.UserManager, itemManager *managers.ItemManager) resources.UserServiceHandler {
	return resources.UserServiceHandler{
		UserSessionManager: userSessionManager,
		UserManager:        userManager,
		ItemManager:        itemManager,
		UserRetriever:      &retrievers.ShelterRetriever{},
	}
}

func buildItemServiceHandler(userSessionManager *managers.UserSessionManager, itemManager *managers.ItemManager, environment *EnvironmentConfig) resources.ItemServiceHandler {
	return resources.ItemServiceHandler{
		UserSessionManager: userSessionManager,
		ItemManager:        itemManager,
		EmailSender:        environment.EmailSender,
		ItemRetriever:      &retrievers.ItemRetriever{},
	}
}

func buildLoginServiceHandler(userSessionManager *managers.UserSessionManager, userManager *managers.UserManager) resources.LoginServiceHandler {
	return resources.LoginServiceHandler{
		UserManager:        userManager,
		UserSessionManager: userSessionManager,
		LoginRetriever:     &retrievers.LoginRetriever{},
	}
}
