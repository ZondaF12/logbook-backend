package profile

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
	store     types.ProfileStore
	userStore types.UserStore
}

func NewHandler(store types.ProfileStore, userStore types.UserStore) *Handler {
	return &Handler{
		store:     store,
		userStore: userStore,
	}
}

func (h *Handler) RegisterRoutes(router *echo.Group) {
	router.POST("/self", auth.WithJWTAuth(h.HandlerCreateProfile, h.userStore))
	router.PUT("/self", auth.WithJWTAuth(h.HandleUpdateProfile, h.userStore))
	router.GET("/self", auth.WithJWTAuth(h.HandleGetProfile, h.userStore))
}

func (h *Handler) HandlerCreateProfile(c echo.Context) error {
	// Parse payload
	var payload types.CreateProfilePayload
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

	// Check if user already has a profile
	_, err := h.store.GetProfileByUserId(userId)
	if err == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Profile already created")
	}

	// Create user
	err = h.store.CreateProfile(types.Profile{
		UserID:   userId,
		Username: payload.Username,
		Name:     payload.Name,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, "User Created")
}

func (h *Handler) HandleUpdateProfile(c echo.Context) error {
	return c.JSON(200, "Update")
}

func (h *Handler) HandleGetProfile(c echo.Context) error {
	userId := auth.GetUserIDFromContext(c.Request().Context())

	u, err := h.store.GetProfileByUserId(userId)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, u)
}
