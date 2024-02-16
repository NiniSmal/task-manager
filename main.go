package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"gitlab.com/nina8884807/task-manager/api"
	"gitlab.com/nina8884807/task-manager/repository"
	"gitlab.com/nina8884807/task-manager/service"
	"log"
	"net/http"
)

func main() {
	//docker run -d -p 8014:5432 -e POSTGRES_PASSWORD=dev -e POSTGRES_DATABASE=postgres postgres

	db, err := sql.Open("postgres", "postgres://postgres:dev@localhost:8014/postgres?sslmode=disable")
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
	err = http.ListenAndServe(":8021", router)
	if err != nil {
		log.Fatal(err)
	}
}
