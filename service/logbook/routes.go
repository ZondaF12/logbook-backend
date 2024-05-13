package logbook

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
	store       types.LogbookStore
	userStore   types.UserStore
	garageStore types.GarageStore
	mediaStore  types.MediaStore
}

func NewHandler(store types.LogbookStore, userStore types.UserStore, garageStore types.GarageStore, mediaStore types.MediaStore) *Handler {
	return &Handler{
		store:       store,
		userStore:   userStore,
		garageStore: garageStore,
		mediaStore:  mediaStore,
	}
}

func (h *Handler) RegisterRoutes(router *echo.Group) {
	router.POST("/log", auth.WithJWTAuth(h.HandleCreateLog, h.userStore))
	router.GET("/log/:vehicleId", auth.WithJWTAuth(h.HandleGetVehicleLogs, h.userStore))
	router.POST("/log/:logId/media", auth.WithJWTAuth(h.HandleUploadLogMedia, h.userStore))
}

func (h *Handler) HandleCreateLog(c echo.Context) error {
	// Parse payload
	var payload types.CreateLogPayload
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

	// Get vehicle from database
	vehicle, err := h.garageStore.GetVehicleByID(payload.VehicleId)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Check if user owns the vehicle
	if vehicle.UserID != userId {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("user does not own vehicle"))
	}

	// Create log
	logId, err := h.store.CreateLog(payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"log_id": logId.String()})
}

func (h *Handler) HandleGetVehicleLogs(c echo.Context) error {
	// Get user ID from JWT
	userId := auth.GetUserIDFromContext(c.Request().Context())

	// Get vehicle ID
	id := c.Param("vehicleId")
	vehicleId := uuid.MustParse(id)

	// Get vehicle from database
	vehicle, err := h.garageStore.GetVehicleByID(vehicleId)
	if err != nil {
		log.Printf("error: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// Check if user owns the vehicle
	if vehicle.UserID != userId {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("user does not own vehicle"))
	}

	// Get logs from database
	logs, err := h.store.GetLogsByVehicleId(vehicleId)
	if err != nil {
		log.Printf("error: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, logs)
}

func (h *Handler) HandleUploadLogMedia(c echo.Context) error {
	id := c.Param("logId")
	logbookId := uuid.MustParse(id)

	// Get avatar file
	file, err := c.FormFile("media")
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
		return c.JSON(http.StatusInternalServerError, "Error uploading media")
	}

	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("logbook-app"),
		Key:    aws.String(fmt.Sprintf("logbook/%d/media/%s", logbookId, file.Filename)),
		Body:   src,
	})
	if err != nil {
		log.Printf("error: %v", err)
		return c.JSON(http.StatusInternalServerError, "Error uploading image")
	}

	fileType := file.Header.Get("Content-Type")
	media := types.Media{
		Filename:   &file.Filename,
		FileType:   &fileType,
		S3Location: &result.Location,
		LogID:      &logbookId,
	}

	// Add media to database
	err = h.mediaStore.AddNewLogMedia(media)
	if err != nil {
		log.Printf("error: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, result.Location)
}
