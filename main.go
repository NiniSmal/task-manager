package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"gitlab.com/nina8884807/task-manager/api"
	"gitlab.com/nina8884807/task-manager/config"
	"gitlab.com/nina8884807/task-manager/middleware"
	"gitlab.com/nina8884807/task-manager/repository"
	"gitlab.com/nina8884807/task-manager/service"
	"log"
	"net/http"
)

func main() {
	ctg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", ctg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db - OK")

	//make migrations

	log.Printf("Start migrating database \n")
	err = goose.Up(db, "migrations")
	if err != nil {
		log.Fatal(err)
	}

	r := repository.NewTaskRepository(db)
	s := service.NewTaskService(r)
	h := api.NewHandler(s)
	router := chi.NewRouter()

	router.Use(middleware.Logging, middleware.ResponseHeader)

	router.HandleFunc("/createTask", h.CreateTask)
	router.HandleFunc("/getTaskByID", h.GetTaskByID)
	router.HandleFunc("/getAllTasks", h.GetAllTasks)
	router.HandleFunc("/updateTask", h.UpdateTask)

	err = http.ListenAndServe(ctg.Port, router)
	if err != nil {
		log.Fatal(err)
	}
}
