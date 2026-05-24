package main

import (
	"context"
	"log"
	"net/http"
	_ "test_em/docs"
	"test_em/internal/config"
	"test_em/internal/handler"
	"test_em/internal/repository"
	"test_em/internal/usecase"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

// @title           Online Subscriptions Aggregator API
// @version         1.0
// @description     REST service for tracking online subscriptions.
// @host            localhost:8080
// @BasePath        /
func main() {
	cfg := config.MustLoad()
	err := repository.RunMigrations(cfg.DBDSN())
	if err != nil && err.Error() != "no change" {
		log.Fatalf("Failed migrations: %v", err)
	}
	log.Println("Migrations successfully applied")
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, cfg.DBDSN())
	if err != nil {
		log.Fatalf("Failed on work of server: %v", err)
	}
	repo := repository.NewPostgresRepository(conn)
	service := usecase.NewService(repo)
	subHandler := handler.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/subscription/create", subHandler.CreateSubscription)
	r.Delete("/subscription/delete", subHandler.DeleteSubscription)
	r.Put("/subscription/update", subHandler.UpdateSubscription)
	r.Get("/subscription", subHandler.GetSubscription)
	r.Get("/subscription/list", subHandler.ListSubscription)
	r.Get("/subscription/total", subHandler.GetTotalCostSub)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Println("Server successfully running on localhost http://localhost:8080")

	err = http.ListenAndServe(cfg.Server.ServerPort, r)
	if err != nil {
		log.Fatalf("Failed on work of server: %v", err)
	}
}
