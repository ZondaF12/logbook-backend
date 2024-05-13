package profile

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ZondaF12/logbook-backend/service/auth"
	"github.com/ZondaF12/logbook-backend/types"
	"github.com/ZondaF12/logbook-backend/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
	router.POST("/self/avatar", auth.WithJWTAuth(h.HandleUploadAvatar, h.userStore))
	router.GET("/user/:id", auth.WithJWTAuth(h.HandleGetUserById, h.userStore))
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

func (h *Handler) HandleUploadAvatar(c echo.Context) error {
	// Get user ID from JWT
	userId := auth.GetUserIDFromContext(c.Request().Context())

	// Get avatar file
	file, err := c.FormFile("avatar")
	if err != nil {
		fmt.Println(err)
		return err
	}
	src, err := file.Open()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer src.Close()

	// Upload avatar to S3
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("error: %v", err)
		return c.JSON(http.StatusInternalServerError, "Error uploading avatar")
	}

	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("logbook-app"),
		Key:    aws.String(fmt.Sprintf("avatars/user/%d/%s", userId, file.Filename)),
		Body:   src,
	})
	if err != nil {
		log.Printf("error: %v", err)
		return c.JSON(http.StatusInternalServerError, "Error uploading avatar")
	}

	// Update user avatar
	err = h.store.UpdateAvatar(userId, result.Location)
	if err != nil {
		log.Printf("error: %v", err)
		return c.JSON(http.StatusInternalServerError, "Error setting avatar in db")
	}

	return c.JSON(http.StatusOK, result.Location)
}

func (h *Handler) HandleGetUserById(c echo.Context) error {
	userId := c.Param("id")
	userIdInt := uuid.MustParse(userId)

	u, err := h.store.GetProfileByUserId(userIdInt)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, u)
}
