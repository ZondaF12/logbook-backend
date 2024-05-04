package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ZondaF12/logbook-backend/config"
	"github.com/ZondaF12/logbook-backend/types"
)

func FetchVehicleDetails(registration string) (types.VehicleInfoRequestData, error) {
	vehicleData, err := DoVehicleInfoRequest(registration)
	if err != nil {
		return types.VehicleInfoRequestData{}, err
	}

	motData, err := DoVehicleMotRequest(registration)
	if err != nil {
		return types.VehicleInfoRequestData{}, err
	}

	var taxDate string
	if vehicleData.TaxDueDate != "" {
		taxDate = vehicleData.TaxDueDate
	} else {
		taxDate = vehicleData.TaxStatus
	}

	var motDate string
	if vehicleData.MotExpiryDate != "" {
		motDate = vehicleData.MotExpiryDate
	} else {
		motDate = motData[0].MotTestExpiryDate
	}

	var registeredDate string
	if motData[0].FirstUsedDate != "" {
		registeredDate = strings.Replace(motData[0].FirstUsedDate, ".", "-", 2)
	} else {
		date, err := time.Parse("2006-01-02", motDate)
		if err != nil {
			fmt.Println(err)
		}

		registeredDate = date.AddDate(-3, 0, 1).Format("2006-01-02")
	}

	newVehicle := types.VehicleInfoRequestData{
		Registration: registration,
		Color:        motData[0].PrimaryColour,
		EngineSize:   uint16(vehicleData.EngineCapacity),
		Make:         vehicleData.Make,
		Model:        motData[0].Model,
		TaxDate:      taxDate,
		MotDate:      motDate,
		Registered:   registeredDate,
		Year:         uint16(vehicleData.YearOfManufacture),
	}

	return newVehicle, nil
}

func DoVehicleInfoRequest(registration string) (types.VehicleData, error) {
	jsonBody := []byte(fmt.Sprintf(`{"registrationNumber": "%s"}`, registration))
	bodyReader := bytes.NewBuffer(jsonBody)

	requestURL := "https://driver-vehicle-licensing.api.gov.uk/vehicle-enquiry/v1/vehicles"
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		return types.VehicleData{}, fmt.Errorf("could not create request %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", config.Envs.DVLAApiKey)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return types.VehicleData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.VehicleData{}, errors.New("invalid registration number")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.VehicleData{}, err
	}

	var vehicleResponse types.VehicleData
	if err = json.Unmarshal(body, &vehicleResponse); err != nil {
		fmt.Printf("Error: %v", err)
	}

	return vehicleResponse, nil
}

func DoVehicleMotRequest(registration string) (types.MotData, error) {
	requestURL := "https://beta.check-mot.service.gov.uk/trade/vehicles/mot-tests/?registration=" + registration
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return types.MotData{}, fmt.Errorf("could not create request %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", config.Envs.DVSAApiKey)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return types.MotData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.MotData{}, errors.New("invalid registration number")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.MotData{}, err
	}

	var vehicleResponse types.MotData
	if err = json.Unmarshal(body, &vehicleResponse); err != nil {
		fmt.Printf("Error: %v", err)
	}

	return vehicleResponse, nil
}
