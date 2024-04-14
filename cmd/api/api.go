package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/ZondaF12/logbook-backend/server/user"
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

	userHandler := user.NewHandler()
	userHandler.RegisterRoutes(subrouter)

	log.Println("Starting server on", s.addr)

	return http.ListenAndServe(s.addr, e)
}
