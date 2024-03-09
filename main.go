package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"gitlab.com/nina8884807/task-manager/api"
	"gitlab.com/nina8884807/task-manager/config"
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
	//подключение к бд
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
	//применяем все возможные миграции
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
	//midll такой же обработчик, поэтому так же принимает репозиторий
	mw := api.NewMiddleware(rt)
	router := chi.NewRouter()

	router.Use(api.Logging, api.ResponseHeader)
	//для части обработчиков создаем группу с доп. middleware для авторизации, тк она нужна не для всех обработчиков
	router.Group(func(r chi.Router) {
		r.Use(mw.AuthHandler)
		r.HandleFunc("/createTask", ht.CreateTask)
		r.HandleFunc("/getTaskByID", ht.GetTaskByID)
		r.HandleFunc("/getAllTasks", ht.GetAllTasks)
	})

	router.HandleFunc("/updateTask", ht.UpdateTask)

	router.HandleFunc("/createUser", hu.CreateUser)
	router.HandleFunc("/login", hu.Login)

	err = http.ListenAndServe(ctg.Port, router)
	if err != nil {
		log.Fatal(err)
	}
}
