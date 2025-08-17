package appServer

import (
	"context"
	"crypto/tls"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ds124wfegd/tech_wildberries_Go/config"
	"github.com/ds124wfegd/tech_wildberries_Go/internal/database"
	"github.com/ds124wfegd/tech_wildberries_Go/internal/service"
	"github.com/ds124wfegd/tech_wildberries_Go/internal/transport"
	"github.com/sirupsen/logrus"
)

type Server struct {
	httpServer *http.Server
}

/*
	AppVersion   string `json:"appVersion"`
	Host         string `json:"host" validate:"required"`
	Port         string `json:"port" validate:"required"`
	Timeout      time.Duration
	Idle_timeout time.Duration
	Env          string `json:"environment"`
*/

func (s *Server) Run(cfg *config.Config, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           handler,
		MaxHeaderBytes:    1 << 20,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		TLSConfig:         &tls.Config{MinVersion: tls.VersionTLS12}, // запрет на устаревшие сертификаты TLS
		//	ErrorLog:          log.New(os.Stderr, "SERVER ERROR: ", log.LstdFlags), //далее можно заменить на ElasticSearch вместо os.StdErr!!!!!
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func NewServer(cfg *config.Config) {

	logrus.SetFormatter(new(logrus.JSONFormatter)) // JSON format for logging

	/* if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}*/

	db, err := database.NewPostgresDB(&config.PostgresConfig{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.DBName,
		SSLMode:  cfg.Postgres.SSLMode,
		PgDriver: cfg.Postgres.PgDriver,
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := database.NewRepository(db)
	services := service.NewService(repos)
	handlers := transport.NewHandler(services)

	srv := new(Server)
	go func() {
		if err := srv.Run(cfg, handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("MpApp Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("MpApp Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	/*if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}*/
}

/*// Инициализация PostgreSQL
pgClient, err := postgres.New(cfg.Postgres.URL)
if err != nil {
	log.Fatalf("failed to connect to postgres: %v", err)
}
defer pgClient.Close()

// Инициализация репозиториев
orderRepo := postgres.NewOrderRepository(pgClient)
cache := in_memory.NewCache()

// Инициализация usecase
orderUC := usecase.NewOrderUseCase(orderRepo, cache)

// Восстановление кэша из БД
if err := orderUC.RestoreCache(context.Background()); err != nil {
	log.Printf("failed to restore cache: %v", err)
}

// Инициализация Kafka Consumer
kafkaConsumer := kafka.NewConsumer(
	cfg.Kafka.Brokers,
	cfg.Kafka.Topic,
	cfg.Kafka.GroupID,
	orderUC,
)
defer kafkaConsumer.Close()
kafkaConsumer.Start(context.Background())
*/

/*	// Инициализация HTTP сервера
	router := gin.Default()
	router.LoadHTMLGlob("internal/delivery/web/template/*")

	// Статические файлы
	router.Static("/static", "./internal/delivery/web/static")

	// API
	orderHandler := http.NewOrderHandler(orderUC)
	router.GET("/order/:order_uid", orderHandler.GetOrder)
	router.GET("/order/html/:order_uid", orderHandler.GetOrderHTML)
	router.GET("/", func(c *gin.Context) {
		c.File("./internal/delivery/web/static/index.html")
	})

	// Запуск сервера
	go func() {
		if err := router.Run(cfg.HTTP.Port); err != nil {
			log.Fatalf("failed to run http server: %v", err)
		}
	}()

	log.Printf("Service started on %s", cfg.HTTP.Port)

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
*/
