package vehicle

import (
	"fmt"
	"net/http"

	"github.com/ZondaF12/logbook-backend/service/auth"
	"github.com/ZondaF12/logbook-backend/types"
	"github.com/ZondaF12/logbook-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	store     types.VehicleStore
	userStore types.UserStore
}

func NewHandler(store types.VehicleStore, userStore types.UserStore) *Handler {
	return &Handler{
		store:     store,
		userStore: userStore,
	}
}

func (h *Handler) RegisterRoutes(router *echo.Group) {
	router.POST("/vehicle", auth.WithJWTAuth(h.HandleAddVehicle, h.userStore))
	router.GET("/vehicle", auth.WithJWTAuth(h.HandleGetAllVehicles, h.userStore))
	router.GET("/vehicle/:registration", auth.WithJWTAuth(h.HandleGetVehicleByRegistration, h.userStore))
	router.PUT("/vehicle/:registration", auth.WithJWTAuth(h.HandleUpdateVehicle, h.userStore))
	router.GET("/vehicle/:registration/exists", auth.WithJWTAuth(h.HandleGetVehicleExists, h.userStore))
}

func (h *Handler) HandleAddVehicle(c echo.Context) error {
	// Parse payload
	var payload types.NewVehiclePostData
	if err := utils.ParseJSON(c, &payload); err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		fmt.Println(err)
		errors := err.(validator.ValidationErrors)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
	}

	// Get user ID from JWT
	userId := auth.GetUserIDFromContext(c.Request().Context())

	// Check if vehicle is already added
	exists, err := h.store.CheckVehicleAdded(userId, payload.Registration)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if exists {
		return echo.NewHTTPError(http.StatusBadRequest, "Vehicle already added")
	}

	// Create vehicle
	err = h.store.AddUserVehicle(userId, payload)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, "Vehicle Added")
}

func (h *Handler) HandleGetAllVehicles(c echo.Context) error {
	// Get user ID from JWT
	userId := auth.GetUserIDFromContext(c.Request().Context())

	vehicles, err := h.store.GetAuthenticatedUserVehicles(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, vehicles)
}

func (h *Handler) HandleGetVehicleByRegistration(c echo.Context) error {
	registration := c.Param("registration")

	// Get user ID from JWT
	userId := auth.GetUserIDFromContext(c.Request().Context())

	vehicle, err := h.store.GetVehicleByRegistration(userId, registration)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, vehicle)
}

func (h *Handler) HandleUpdateVehicle(c echo.Context) error {
	return nil
}

func (h *Handler) HandleGetVehicleExists(c echo.Context) error {
	registration := c.Param("registration")

	// Get user ID from JWT
	userId := auth.GetUserIDFromContext(c.Request().Context())

	exists, err := h.store.CheckVehicleAdded(userId, registration)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, exists)
}
