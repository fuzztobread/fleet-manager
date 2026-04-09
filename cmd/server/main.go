// cmd/server/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"fleet-manager/internal/dispatch"
	"fleet-manager/internal/route"
	"fleet-manager/internal/storage"
	"fleet-manager/internal/vehicle"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	db, err := storage.New("fleet.db")
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}
	defer db.Close()

	vehicleStore := storage.NewVehicleStore(db)
	vehicleService := vehicle.NewService(vehicleStore)
	routeService := route.NewService()
	dispatchService := dispatch.NewService(vehicleStore, routeService)

	// context that cancels on SIGINT/SIGTERM — propagates to the worker
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// background worker — drains the priority queue
	go dispatchService.Run(ctx)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/vehicles", vehicle.NewHandler(vehicleService))
	r.Mount("/routes", route.NewHandler(routeService))
	r.Mount("/dispatch", dispatch.NewHandler(dispatchService))

	log.Println("starting fleet-manager on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
