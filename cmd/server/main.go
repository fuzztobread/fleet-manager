package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"fleet-manager/internal/storage"
	"fleet-manager/internal/vehicle"
)

func main() {
	db, err := storage.New("fleet.db")
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}
	defer db.Close() // runs when main() returns (e.g. on shutdown signal)

	vehicleStore := storage.NewVehicleStore(db)

	vehicleService := vehicle.NewService(vehicleStore)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/vehicles", vehicle.NewHandler(vehicleService))

	log.Println("starting fleet-manager on :8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
