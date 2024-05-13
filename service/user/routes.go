package user

import (
	"fmt"
	"net/http"

	"github.com/ZondaF12/logbook-backend/config"
	"github.com/ZondaF12/logbook-backend/service/auth"
	"github.com/ZondaF12/logbook-backend/types"
	"github.com/ZondaF12/logbook-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *echo.Group) {
	router.POST("/login", h.HandleLogin)
	router.POST("/register", h.HandleRegister)
}

func (h *Handler) HandleLogin(c echo.Context) error {
	// Parse payload
	var payload types.LoginAuthPayload
	if err := utils.ParseJSON(c, &payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
	}

	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
	}

	if !auth.ComparePassword(u.Password, payload.Password) {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.Response().Header().Set("Access-Control-Expose-Headers", "X-Logbook-Token")
	c.Response().Header().Set("X-Logbook-Token", token)

	return c.JSON(http.StatusOK, map[string]string{"userId": u.ID.String()})
}

func (h *Handler) HandleRegister(c echo.Context) error {
	// Parse payload
	var payload types.RegisterAuthPayload
	if err := utils.ParseJSON(c, &payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
	}

	// Check if user exists
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// Create user
	err = h.store.CreateUser(types.User{
		Email:    payload.Email,
		Password: hashedPassword,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, "User Created")
}
