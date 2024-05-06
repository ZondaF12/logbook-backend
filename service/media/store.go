package media

import (
	"database/sql"

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

func (s *Store) AddNewVehicleMedia(media types.Media) error {
	_, err := s.db.Exec(`
		INSERT INTO media (filename, file_type, s3_location, vehicle_id)
		VALUES (?, ?, ?, ?)`,
		media.Filename, media.FileType, media.S3Location, media.VehicleID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) AddNewLogMedia(media types.Media) error {
	_, err := s.db.Exec(`
		INSERT INTO media (filename, file_type, s3_location, log_id)
		VALUES (?, ?, ?, ?)`,
		media.Filename, media.FileType, media.S3Location, media.LogID,
	)
	if err != nil {
		return err
	}

	return nil
}
