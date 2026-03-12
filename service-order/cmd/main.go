package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	orderHandler "github.com/barashF/lms/service-order/internal/handler/order"
	"github.com/barashF/lms/service-order/internal/logger"
	"github.com/barashF/lms/service-order/internal/publisher"
	orderRepo "github.com/barashF/lms/service-order/internal/repository/order"
	"github.com/barashF/lms/service-order/internal/repository/outbox"
	"github.com/barashF/lms/service-order/internal/repository/utils/transaction"
	orderServ "github.com/barashF/lms/service-order/internal/service/order"
	"github.com/barashF/lms/service-order/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appLogger, err := logger.NewZapAdapter()
	if err != nil {
		log.Fatalf("failed initialize logger: %v", err)
	}
	defer appLogger.Sync()

	err = godotenv.Load()
	if err != nil {
		appLogger.Error("Error loading .env file", logger.NewField("error", err.Error()))
	}

	dbPool := database.MustInitDB()
	transactionManager := transaction.NewManager(dbPool)
	orderRepository := orderRepo.NewRepository(transactionManager)
	outboxRepository := outbox.NewRepository(transactionManager)
	orderService := orderServ.NewService(orderRepository, outboxRepository, transactionManager, appLogger)

	kafkaPublisher, err := publisher.NewKafkaPublisher([]string{"kafka-1:9092", "kafka-2:9092", "kafka-3:9092"})
	if err != nil {
		log.Fatalf("Failed to create Kafka publisher: %v", err)
	}
	defer kafkaPublisher.Close()
	workerPublicher := publisher.NewPublisher(outboxRepository, kafkaPublisher, appLogger, time.Second*10, 5)
	go workerPublicher.Start(ctx)

	startServer(
		dbPool,
		orderHandler.NewController(orderService),
		resolvePort(),
		appLogger,
	)
}

func startServer(
	dbPool *pgxpool.Pool,
	orderHandler *orderHandler.Controller,
	port string,
	appLogger logger.Logger,
) {
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: initRouter(orderHandler),
	}
	appLogger.Info("Server started on ", logger.NewField("address", srv.Addr))

	serverErr := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	waitGracefullShutdown(srv, dbPool, serverErr, appLogger)
	appLogger.Info("Shutting down service-order")
}

func resolvePort() string {
	port := os.Getenv("PORT")

	portFlag := flag.String("port", "", "укажите порт")
	flag.Parse()

	if portFlag != nil && *portFlag != "" {
		port = *portFlag
	}

	if port == "" {
		port = "8080"
	}

	return port
}

func initRouter(order *orderHandler.Controller) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/order", func(r chi.Router) {
		r.Post("/", order.Create)
	})
	return r
}

func waitGracefullShutdown(srv *http.Server, dbPool *pgxpool.Pool, serverErr <-chan error, appLogger logger.Logger) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var reason string
	select {
	case <-ctx.Done():
		reason = "signal"
	case err := <-serverErr:
		reason = "server error: " + err.Error()
	}
	appLogger.Info("Shutdown initiated", logger.NewField("reason", reason))

	appLogger.Info("Shutting down HTTP server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		appLogger.Error("HTTP server shutdown", logger.NewField("error", err))
	} else {
		appLogger.Info("HTTP server stopped")
	}

	appLogger.Info("Closing DB pool...")
	dbPool.Close()
	appLogger.Info("DB pool closed")
}
