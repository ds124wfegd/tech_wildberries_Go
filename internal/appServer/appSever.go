// launching the server, DB, kafka, postgres
package appServer

import (
	"context"
	"crypto/tls"
	"log"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ds124wfegd/tech_wildberries_Go/config"
	"github.com/ds124wfegd/tech_wildberries_Go/internal/cache"
	"github.com/ds124wfegd/tech_wildberries_Go/internal/database"
	"github.com/ds124wfegd/tech_wildberries_Go/internal/entity"
	"github.com/ds124wfegd/tech_wildberries_Go/internal/kafka"
	"github.com/ds124wfegd/tech_wildberries_Go/internal/service"
	"github.com/ds124wfegd/tech_wildberries_Go/internal/transport"

	"github.com/sirupsen/logrus"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(cfg *config.Config, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           handler,
		MaxHeaderBytes:    1 << 20,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      cfg.Server.Timeout,
		IdleTimeout:       cfg.Server.Idle_timeout,
		ReadHeaderTimeout: 3 * time.Second,
		TLSConfig:         &tls.Config{MinVersion: tls.VersionTLS12},           // ban on outdate TLS certificate
		ErrorLog:          log.New(os.Stderr, "SERVER ERROR: ", log.LstdFlags), // os.Stderr can be replaced with ElsasticSearch in the feature
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func NewServer(cfg *config.Config) {

	logrus.SetFormatter(new(logrus.JSONFormatter)) // JSON format for logging

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

	orderCache := cache.NewCache()

	var repo database.OrderRepository
	if db != nil {
		repo = repos
	}

	var cachePort service.OrderCache = orderCache

	services := service.NewService(repos, cachePort)
	handlers := transport.NewOrderHandler(services)

	// restoring the cache from the db
	if repo != nil {
		logrus.Println("Restoring cache from database...")
		uids, err := repo.GetRecentUIDs(context.Background(), 100)
		if err != nil {
			logrus.Printf("Warning: Failed to fetch recent order uids: %v", err)
		} else {
			var restoreOrders []*entity.Order
			for _, uid := range uids {
				o, err := repo.GetByUID(context.Background(), uid)
				if err != nil {
					logrus.Printf("Warning: Failed to get order %s: %v", uid, err)
					continue
				}
				restoreOrders = append(restoreOrders, o)
			}
			cachePort.Load(restoreOrders)
			logrus.Printf("Restored %d orders to cache", len(restoreOrders))
		}
	}

	kafkaConsumer := kafka.NewConsumer(
		[]string{config.GetEnv("KAFKA_BROKERS", "localhost:9094")},
		config.GetEnv("KAFKA_TOPIC", "orders"),
		config.GetEnv("KAFKA_GROUP_ID", "order-service"),
		services,
	)

	// launching Kafka consumer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { kafkaConsumer.Start(ctx) }()

	srv := new(Server)
	go func() {
		if err := srv.Run(cfg, handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("App Started")

	// wating for the completion signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("App Shutting Down")

	//shutdown - correct completion of the program
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}

}
