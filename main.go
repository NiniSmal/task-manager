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
	gen "gitlab.com/nina8884807/mail/proto"
	"gitlab.com/nina8884807/task-manager/api"
	"gitlab.com/nina8884807/task-manager/config"
	"gitlab.com/nina8884807/task-manager/repository"
	"gitlab.com/nina8884807/task-manager/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	defer db.Close()

	ctx := context.Background()
	rds := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
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

	con, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Errorf("dial: %w", err)
	}
	mailClient := gen.NewMailClient(con)

	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "topic-A",
		Balancer: &kafka.LeastBytes{},
	}
	defer kafkaWriter.Close()

	rt := repository.NewTaskRepository(db, rds)
	st := service.NewTaskService(rt)
	ht := api.NewTaskHandler(st)
	ut := repository.NewUserRepository(db, rds)
	su := service.NewUserService(ut, mailClient, kafkaWriter)
	hu := api.NewUserHandler(su)
	//midll такой же обработчик, поэтому так же принимает репозиторий
	mw := api.NewMiddleware(ut)
	router := chi.NewRouter()

	router.Use(api.Logging, api.ResponseHeader)
	//для части обработчиков создаем группу с доп. middleware для авторизации, тк она нужна не для всех обработчиков
	router.Group(func(r chi.Router) {
		r.Use(mw.AuthHandler)
		r.Post("/tasks", ht.CreateTask)
		r.Get("/tasks/{id}", ht.GetTaskByID)
		r.Get("/tasks", ht.GetAllTasks)
		r.Put("/tasks/{id}", ht.UpdateTask)
	})

	router.Post("/users", hu.CreateUser)
	router.Post("/login", hu.Login)
	router.Get("/verification", hu.Verification)

	err = http.ListenAndServe(ctg.Port, router)
	if err != nil {
		log.Fatal(err)
	}
}
