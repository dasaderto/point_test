package main

import (
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
	"news/db"
	"news/repositories"
	"os"
)

func main() {
	initLogger()
	loadEnv()
	baseRepos := repositories.NewRepository(db.DbConnect())
	baseRepos.Migrate()
	router := MakeRouter()
	http.Handle("/", router)
	fmt.Println(fmt.Sprintf("Listen port %s", os.Getenv("LISTEN_PORT")))
	err := http.ListenAndServe(os.Getenv("LISTEN_PORT"), nil)
	if err != nil{
		log.Error(err.Error())
	}
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func initLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}
