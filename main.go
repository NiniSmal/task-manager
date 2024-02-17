package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"gitlab.com/nina8884807/task-manager/api"
	"gitlab.com/nina8884807/task-manager/config"
	"gitlab.com/nina8884807/task-manager/repository"
	"gitlab.com/nina8884807/task-manager/service"
	"log"
	"net/http"
)

func main() {
	ctg := config.GetConfig()

	db, err := sql.Open("postgres", ctg.Data)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db - OK")

	r := repository.NewTaskRepository(db)
	s := service.NewTaskService(r)
	h := api.NewHandler(s)
	router := http.NewServeMux()
	router.HandleFunc("/createTask", h.CreateTask)
	err = http.ListenAndServe(ctg.Port, router)
	if err != nil {
		log.Fatal(err)
	}
}
