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

	rt := repository.NewTaskRepository(db)
	st := service.NewTaskService(rt)
	ht := api.NewTaskHandler(st)
	ut := repository.NewUserRepository(db)
	su := service.NewUserService(ut)
	hu := api.NewUserHandler(su)

	router := chi.NewRouter()

	router.Use(middleware.Logging, middleware.ResponseHeader)

	router.HandleFunc("/createTask", ht.CreateTask)
	router.HandleFunc("/getTaskByID", ht.GetTaskByID)
	router.HandleFunc("/getAllTasks", ht.GetAllTasks)
	router.HandleFunc("/updateTask", ht.UpdateTask)

	router.HandleFunc("/createUser", hu.CreateUser)
	router.HandleFunc("/login", hu.Login)

	err = http.ListenAndServe(ctg.Port, router)
	if err != nil {
		log.Fatal(err)
	}
}
