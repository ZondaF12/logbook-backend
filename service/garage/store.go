package garage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ZondaF12/logbook-backend/types"
	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetAuthenticatedUserVehicles(userId uuid.UUID) ([]*types.Vehicle, error) {
	rows, err := s.db.Query(`SELECT
			v.*,
			(
			SELECT
				GROUP_CONCAT(m.s3_location)
			FROM
				media m
			WHERE
				m.vehicle_id = v.id
			) AS media
		FROM
			vehicles v
		WHERE
			v.user_id = ?
		ORDER BY
			v.created_at`, userId)
	if err != nil {
		return nil, err
	}

	vehicles := make([]*types.Vehicle, 0)
	for rows.Next() {
		vehicle, err := scanRowIntoVehicle(rows)
		if err != nil {
			return nil, err
		}

		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}

func scanRowIntoVehicle(rows *sql.Rows) (*types.Vehicle, error) {
	vehicle := new(types.Vehicle)

	err := rows.Scan(
		&vehicle.ID,
		&vehicle.UserID,
		&vehicle.Registration,
		&vehicle.Make,
		&vehicle.Model,
		&vehicle.Year,
		&vehicle.EngineSize,
		&vehicle.Color,
		&vehicle.Registered,
		&vehicle.TaxDate,
		&vehicle.MotDate,
		&vehicle.InsuranceDate,
		&vehicle.ServiceDate,
		&vehicle.Description,
		&vehicle.Mileage,
		&vehicle.Nickname,
		&vehicle.CreatedAt,
		&vehicle.Images,
	)
	if err != nil {
		return nil, err
	}

	return vehicle, nil
}

func (s *Store) GetVehicleByRegistration(userId uuid.UUID, registration string) (*types.Vehicle, error) {
	rows, err := s.db.Query("SELECT * FROM vehicles WHERE user_id = ? AND registration = ?", userId, registration)
	if err != nil {
		return nil, err
	}

	vehicle := new(types.Vehicle)
	for rows.Next() {
		vehicle, err = scanRowIntoVehicle(rows)
		if err != nil {
			return nil, err
		}
	}

	return vehicle, nil
}

func (s *Store) CheckVehicleAdded(userId uuid.UUID, registration string) (bool, error) {
	rows, err := s.db.Query("SELECT EXISTS(SELECT 1 FROM vehicles WHERE user_id = ? AND registration = ?)", userId, registration)
	if err != nil {
		return false, err
	}

	var exists bool
	for rows.Next() {
		err = rows.Scan(&exists)
		if err != nil {
			return false, err
		}
	}

	return exists, nil
}

func (s *Store) GetVehicleByID(vehicleId uuid.UUID) (*types.Vehicle, error) {
	rows, err := s.db.Query(`SELECT
			v.*,
			(
			SELECT
				GROUP_CONCAT(m.s3_location)
			FROM
				media m
			WHERE
				m.vehicle_id = v.id
			) AS media
		FROM
			vehicles v
		WHERE
			v.id = ?
		ORDER BY
			v.created_at`, vehicleId)
	if err != nil {
		return nil, err
	}

	vehicle := new(types.Vehicle)
	for rows.Next() {
		vehicle, err = scanRowIntoVehicle(rows)
		if err != nil {
			return nil, err
		}
	}

	return vehicle, nil
}

func (s *Store) AddUserVehicle(userId uuid.UUID, vehicle types.NewVehiclePostData) (uuid.UUID, error) {
	newVehicleId := uuid.New()
	_, err := s.db.Exec("INSERT INTO vehicles (id, user_id, registration, make, model, year, engine_size, color, registered, tax_date, mot_date, insurance_date, service_date, description, milage, nickname) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", newVehicleId, userId, vehicle.Registration, vehicle.Make, vehicle.Model, vehicle.Year, vehicle.EngineSize, vehicle.Color, vehicle.Registered, vehicle.TaxDate, vehicle.MotDate, "", "", vehicle.Description, vehicle.Mileage, vehicle.Nickname)

	if err != nil {
		fmt.Println(err)
		return uuid.Nil, err
	}

	return newVehicleId, nil
}

func (s *Store) UpdateVehicle(userId uuid.UUID, registration string, vehicle types.UpdateVehiclePatchData) error {
	// Build the SQL query dynamically based on the provided values
	query := "UPDATE vehicles SET"
	params := []interface{}{}

	// Check if other fields are provided and add them to the query and params
	if vehicle.Description != "" {
		query += " description = ?,"
		params = append(params, vehicle.Description)
	}
	if vehicle.MotDate != "" {
		query += " mot_date = ?,"
		params = append(params, vehicle.MotDate)
	}
	if vehicle.InsuranceDate != "" {
		query += " insurance_date = ?,"
		params = append(params, vehicle.InsuranceDate)
	}
	if vehicle.ServiceDate != "" {
		query += " service_date = ?,"
		params = append(params, vehicle.ServiceDate)
	}
	if vehicle.TaxDate != "" {
		query += " tax_date = ?,"
		params = append(params, vehicle.TaxDate)
	}
	if vehicle.Mileage != 0 {
		query += " milage = ?,"
		params = append(params, vehicle.Mileage)
	}
	if vehicle.Nickname != "" {
		query += " nickname = ?,"
		params = append(params, vehicle.Nickname)
	}

	hasTrailingComma := strings.HasSuffix(query, ",")
	if hasTrailingComma {
		query = query[:len(query)-1]
	}

	// Add the WHERE clause to specify the user_id and registration
	query += " WHERE user_id = ? AND registration = ?"
	params = append(params, userId, registration)

	fmt.Println(query)
	fmt.Println(params)

	// Execute the dynamic query with the provided parameters
	_, err := s.db.Exec(query, params...)
	if err != nil {
		return err
	}

	return nil
}
