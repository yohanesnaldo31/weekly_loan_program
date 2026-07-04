package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"weekly_loan_program/app"
)

const port = 8000

func main() {
	application, err := app.NewApplication(context.Background())
	if err != nil {
		log.Fatal("Fail initializing application: ", err)
	}

	application.HTTPHandler.RegisterRoutes()

	server := http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: application.HTTPHandler.Mux,
	}

	log.Println("Starting Server on Port ", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Println("Fail serving: ", err)
	}
}
