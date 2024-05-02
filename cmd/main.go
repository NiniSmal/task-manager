package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"gitlab.com/nina8884807/task-manager/api"
	"gitlab.com/nina8884807/task-manager/config"
	"gitlab.com/nina8884807/task-manager/repository"
	"gitlab.com/nina8884807/task-manager/service"
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

	h := slog.Handler(slog.NewTextHandler(os.Stdout, nil))
	if cfg.LogJson {
		h = slog.NewTextHandler(os.Stdout, nil)
	}
	logger := slog.New(h)

	db, err := sql.Open("postgres", cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Connected to Postgres OK")
	defer db.Close()

	ctx := context.Background()
	rds := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: "",
		DB:       0,
	})

	_, err = rds.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	logger.Info("Connected to Redis OK")

	defer rds.Close()

	// make migrations
	logger.Info("Start migrating database")
	// применяем все возможные миграции
	goose.SetLogger(goose.NopLogger())
	err = goose.Up(db, "migrations")
	if err != nil {
		log.Fatal(err)
	}

	connKafka, err := kafka.Dial("tcp", cfg.KafkaAddr)
	if err != nil {
		log.Fatal("dial kafka:", err)
	}

	connKafka.Close()

	topicConfigs := []kafka.TopicConfig{
		{Topic: cfg.KafkaTopicCreateUser,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}
	err = connKafka.CreateTopics(topicConfigs...)
	if err != nil {
		logger.Error("create topic", err)
	}

	defer connKafka.Close()
	kafkaWriter := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.KafkaAddr),
		Topic:                  cfg.KafkaTopicCreateUser,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	defer kafkaWriter.Close()

	appURL, err := url.Parse(cfg.AppURL)
	if err != nil {
		log.Fatal(err)
	}

	rp := repository.NewProjectRepository(db, rds)
	sp := service.NewProjectService(rp)
	hp := api.NewProjectHandler(sp)
	rt := repository.NewTaskRepository(db, rds)
	st := service.NewTaskService(rt, rp, kafkaWriter)
	ht := api.NewTaskHandler(st)
	ut := repository.NewUserRepository(db, rds)
	su := service.NewUserService(ut, kafkaWriter, cfg.AppURL)
	hu := api.NewUserHandler(su, appURL.Hostname())
	// midll такой же обработчик, поэтому так же принимает репозиторий
	mw := api.NewMiddleware(ut, logger)
	router := chi.NewRouter()

	router.Use(mw.Logging, mw.ResponseHeader)
	// для части обработчиков создаем группу с доп. middleware для авторизации, тк она нужна не для всех обработчиков
	router.Group(func(r chi.Router) {
		r.Use(mw.AuthHandler)
		r.Post("/api/projects", hp.CreateProject)
		r.Get("/api/projects", hp.Projects)
		r.Get("/api/projects/{id}", hp.ProjectByID)
		r.Put("/api/projects/{id}", hp.UpdateProject)
		r.Delete("/api/projects/{id}", hp.DeleteProject)
		r.Post("/api/project/users", hp.AddProjectMember)
		r.Get("/api/user/projects", hp.UserProjects)

		r.Get("/api/tasks", ht.GetAllTasks)
		r.Post("/api/tasks", ht.CreateTask)
		r.Get("/api/tasks/{id}", ht.GetTaskByID)
		r.Put("/api/tasks/{id}", ht.UpdateTask)
	})

	router.Post("/api/users", hu.CreateUser)
	router.Post("/api/login", hu.Login)
	router.Get("/api/verification", hu.Verification)

	logger.Info(fmt.Sprintf("start http server at port: %v", cfg.Port))

	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
	if err != nil {
		log.Fatal(err)
	}
}
