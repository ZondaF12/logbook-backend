package garage

import (
	"database/sql"
	"fmt"

	"github.com/ZondaF12/logbook-backend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetAuthenticatedUserVehicles(userId int) ([]*types.Vehicle, error) {
	rows, err := s.db.Query("SELECT * FROM vehicles WHERE user_id = ?", userId)
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
	)
	if err != nil {
		return nil, err
	}

	return vehicle, nil
}

func (s *Store) GetVehicleByRegistration(userId int, registration string) (*types.Vehicle, error) {
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

func (s *Store) CheckVehicleAdded(userId int, registration string) (bool, error) {
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

func (s *Store) GetVehicleByID(vehicleId int) (*types.Vehicle, error) {
	rows, err := s.db.Query("SELECT * FROM vehicles WHERE id = ?", vehicleId)
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

func (s *Store) AddUserVehicle(userId int, vehicle types.NewVehiclePostData) error {
	_, err := s.db.Exec("INSERT INTO vehicles (user_id, registration, make, model, year, engine_size, color, registered, tax_date, mot_date, insurance_date, service_date, description, milage, nickname) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", userId, vehicle.Registration, vehicle.Make, vehicle.Model, vehicle.Year, vehicle.EngineSize, vehicle.Color, vehicle.Registered, vehicle.TaxDate, vehicle.MotDate, "", "", vehicle.Description, 0, vehicle.Nickname)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
