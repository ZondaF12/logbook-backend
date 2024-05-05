package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/ZondaF12/logbook-backend/service/follower"
	"github.com/ZondaF12/logbook-backend/service/garage"
	"github.com/ZondaF12/logbook-backend/service/image"
	"github.com/ZondaF12/logbook-backend/service/profile"
	"github.com/ZondaF12/logbook-backend/service/user"
	"github.com/ZondaF12/logbook-backend/service/vehicle"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Start() error {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	subrouter := e.Group("/api/v1")

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	profileStore := profile.NewStore(s.db)
	profileHandler := profile.NewHandler(profileStore, userStore)
	profileHandler.RegisterRoutes(subrouter)

	followStore := follower.NewStore(s.db)
	followHandler := follower.NewHandler(followStore, userStore)
	followHandler.RegisterRoutes(subrouter)

	imageStore := image.NewStore(s.db)

	garageStore := garage.NewStore(s.db)
	garageHandler := garage.NewHandler(garageStore, userStore, imageStore)
	garageHandler.RegisterRoutes(subrouter)

	vehicleHandler := vehicle.NewHandler(userStore)
	vehicleHandler.RegisterRoutes(subrouter)

	log.Println("Starting server on", s.addr)
	return http.ListenAndServe(s.addr, e)
}
