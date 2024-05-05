package image

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

func (s *Store) AddNewImage(image types.Image) error {
	_, err := s.db.Exec(`
		INSERT INTO images (filename, file_type, s3_location, user_id, vehicle_id)
		VALUES (?, ?, ?, ?, ?)`,
		image.Filename, image.FileType, image.S3Location, image.UserID, image.VehicleID,
	)
	if err != nil {
		return err
	}

	return nil
}
