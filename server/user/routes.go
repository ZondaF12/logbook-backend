package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *echo.Group) {
	router.POST("/login", h.HandleLogin)
	router.POST("/register", h.HandleRegister)

}

func (h *Handler) HandleLogin(c echo.Context) error {
	return c.JSON(http.StatusOK, "Login")
}

func (h *Handler) HandleRegister(c echo.Context) error {
	return c.JSON(http.StatusCreated, "Register")
}
