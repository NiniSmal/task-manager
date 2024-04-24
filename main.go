package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"gitlab.com/nina8884807/task-manager/api"
	"gitlab.com/nina8884807/task-manager/config"
	"gitlab.com/nina8884807/task-manager/repository"
	"gitlab.com/nina8884807/task-manager/service"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = cfg.Validation()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", cfg)

	db, err := sql.Open("postgres", cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db - OK")
	defer db.Close()

	ctx := context.Background()
	rds := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: "",
		DB:       0,
	})

	pong, err := rds.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	fmt.Println("Connected to Redis:", pong)

	defer rds.Close()

	//make migrations
	log.Printf("Start migrating database \n")
	//применяем все возможные миграции
	err = goose.Up(db, "migrations")
	if err != nil {
		log.Fatal(err)
	}

	connKafka, err := kafka.DialLeader(ctx, "tcp", cfg.KafkaAddr, cfg.KafkaTopicCreateUser, 0)
	if err != nil {
		log.Fatal("dial kafka:", err)
	}

	defer connKafka.Close()

	rp := repository.NewProjectRepository(db, rds)
	sp := service.NewProjectService(rp)
	hp := api.NewProjectHandler(sp)
	rt := repository.NewTaskRepository(db, rds)
	st := service.NewTaskService(rt)
	ht := api.NewTaskHandler(st)
	ut := repository.NewUserRepository(db, rds)
	su := service.NewUserService(ut, connKafka)
	hu := api.NewUserHandler(su)
	//midll такой же обработчик, поэтому так же принимает репозиторий
	mw := api.NewMiddleware(ut)
	router := chi.NewRouter()

	router.Use(api.Logging, api.ResponseHeader)
	//для части обработчиков создаем группу с доп. middleware для авторизации, тк она нужна не для всех обработчиков
	router.Group(func(r chi.Router) {
		r.Use(mw.AuthHandler)
		r.Post("/projects", hp.CreateProject)
		r.Get("/projects", hp.GetAllProjects)
		r.Get("/projects/{id}", hp.GetProject)
		r.Put("/projects/{id}", hp.UpdateProject)
		r.Post("/tasks", ht.CreateTask)
		r.Get("/tasks/{id}", ht.GetTaskByID)
		r.Get("/tasks", ht.GetAllTasks)
		r.Put("/tasks/{id}", ht.UpdateTask)
	})

	router.Post("/users", hu.CreateUser)
	router.Post("/login", hu.Login)
	router.Get("/verification", hu.Verification)

	log.Println("start http server at port:", cfg.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
	if err != nil {
		log.Fatal(err)
	}
}
