package main

import (
	"context"
	"log"
	"time"

	"weekly_loan_program/app"
)

func main() {
	ctx := context.Background()
	application, err := app.NewApplication(ctx)
	if err != nil {
		log.Fatal("fail initializing application: ", err)
	}

	err = application.CronHandler.ImplementTasks()
	if err != nil {
		log.Fatal(err)
	}

	application.CronHandler.Start()
	defer application.CronHandler.Stop()

	log.Println("cron scheduler started")
	<-ctx.Done()
	time.Sleep(1 * time.Second)
}
