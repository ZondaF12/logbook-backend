package vehicle

import (
	"github.com/ZondaF12/logbook-backend/types"
	"github.com/ZondaF12/logbook-backend/utils"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	userStore types.UserStore
}

func NewHandler(userStore types.UserStore) *Handler {
	return &Handler{userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *echo.Group) {
	router.GET("/vehicle/:registration/getDetails", h.GetVehicleDetails)
}

func (h *Handler) GetVehicleDetails(c echo.Context) error {
	registration := c.Param("registration")

	vehicleData, err := utils.FetchVehicleDetails(registration)
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	return c.JSON(200, vehicleData)
}
