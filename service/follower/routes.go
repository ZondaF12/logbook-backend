package follower

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
	store     types.FollowerStore
	userStore types.UserStore
}

func NewHandler(store types.FollowerStore, userStore types.UserStore) *Handler {
	return &Handler{
		store:     store,
		userStore: userStore,
	}
}

func (h *Handler) RegisterRoutes(router *echo.Group) {
	router.POST("/follow", auth.WithJWTAuth(h.HandleFollowUser, h.userStore))
	router.POST("/unfollow", auth.WithJWTAuth(h.HandleUnfollowUser, h.userStore))
}

func (h *Handler) HandleFollowUser(c echo.Context) error {
	// Parse payload
	var payload types.FollowUserPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
	}

	// Get user ID from JWT
	userId := auth.GetUserIDFromContext(c.Request().Context())

	if userId == payload.UserID {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot follow yourself")
	}

	// Check user isnt already following
	f, err := h.store.GetFollower(userId, payload.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if f != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Already following user")
	}

	// Follow user
	err = h.store.FollowUser(userId, payload.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.String(http.StatusOK, fmt.Sprintf("Now following user %d", payload.UserID))
}

func (h *Handler) HandleUnfollowUser(c echo.Context) error {
	// Parse payload
	var payload types.FollowUserPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
	}

	// Get user ID from JWT
	userId := auth.GetUserIDFromContext(c.Request().Context())

	if userId == payload.UserID {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot unfollow yourself")
	}

	// Check user isnt already following
	f, err := h.store.GetFollower(userId, payload.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if f == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Not following user")
	}

	// Unfollow user
	fmt.Println(userId, payload.UserID)
	err = h.store.UnfollowUser(userId, payload.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.String(http.StatusOK, fmt.Sprintf("Unfollowed user %d", payload.UserID))
}
