package garage

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
	store      types.GarageStore
	userStore  types.UserStore
	mediaStore types.MediaStore
}

func NewHandler(store types.GarageStore, userStore types.UserStore, mediaStore types.MediaStore) *Handler {
	return &Handler{
		store:      store,
		userStore:  userStore,
		mediaStore: mediaStore,
	}
}

func (h *Handler) RegisterRoutes(router *echo.Group) {
	router.POST("/garage/vehicle", auth.WithJWTAuth(h.HandleAddVehicleToGarage, h.userStore))
	router.GET("/garage", auth.WithJWTAuth(h.HandleGetUserGarage, h.userStore))
	router.GET("/garage/vehicle/:registration", auth.WithJWTAuth(h.HandleGetVehicleByRegistration, h.userStore))
	router.PATCH("/garage/vehicle/:registration", auth.WithJWTAuth(h.HandleUpdateVehicle, h.userStore))
	router.GET("/garage/vehicle/:registration/exists", auth.WithJWTAuth(h.HandleCheckVehicleExistsInGarage, h.userStore))
	router.POST("/garage/vehicle/:id/uploadImage", auth.WithJWTAuth(h.HandleUploadVehicleImage, h.userStore))
}

func (h *Handler) HandleAddVehicleToGarage(c echo.Context) error {
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
	fmt.Println("Getting user ID from context")
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
	vehicleId, err := h.store.AddUserVehicle(userId, payload)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, map[string]string{"vehicle_id": vehicleId.String()})
}

func (h *Handler) HandleGetUserGarage(c echo.Context) error {
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
	registration := c.Param("registration")

	// Get user ID from JWT
	userId := auth.GetUserIDFromContext(c.Request().Context())

	// Parse payload
	var payload types.UpdateVehiclePatchData
	if err := utils.ParseJSON(c, &payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
	}

	err := h.store.UpdateVehicle(userId, registration, payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "Vehicle updated")
}

func (h *Handler) HandleCheckVehicleExistsInGarage(c echo.Context) error {
	registration := c.Param("registration")

	// Get user ID from JWT
	userId := auth.GetUserIDFromContext(c.Request().Context())

	exists, err := h.store.CheckVehicleAdded(userId, registration)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, exists)
}

func (h *Handler) HandleUploadVehicleImage(c echo.Context) error {
	fmt.Println("Uploading image")
	id := c.Param("id")
	vehicleId := uuid.MustParse(id)

	// Get avatar file
	file, err := c.FormFile("image")
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
		return c.JSON(http.StatusInternalServerError, "Error uploading image")
	}

	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("logbook-app"),
		Key:    aws.String(fmt.Sprintf("vehicles/%d/images/%s", vehicleId, file.Filename)),
		Body:   src,
	})
	if err != nil {
		log.Printf("error: %v", err)
		return c.JSON(http.StatusInternalServerError, "Error uploading image")
	}

	fileType := file.Header.Get("Content-Type")
	userID := auth.GetUserIDFromContext(c.Request().Context())
	media := types.Media{
		Filename:   &file.Filename,
		FileType:   &fileType,
		S3Location: &result.Location,
		VehicleID:  &vehicleId,
		UserID:     &userID,
	}

	// Add media to database
	err = h.mediaStore.AddNewVehicleMedia(media)
	if err != nil {
		log.Printf("error: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, result.Location)
}
