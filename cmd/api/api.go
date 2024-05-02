package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/ZondaF12/logbook-backend/service/profile"
	"github.com/ZondaF12/logbook-backend/service/user"
	"github.com/labstack/echo/v4"
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
	subrouter := e.Group("/api/v1")

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	profileStore := profile.NewStore(s.db)
	profileHandler := profile.NewHandler(profileStore, userStore)
	profileHandler.RegisterRoutes(subrouter)

	log.Println("Starting server on", s.addr)
	return http.ListenAndServe(s.addr, e)
}
