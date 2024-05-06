package logbook

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

func (s *Store) CreateLog(log types.CreateLogPayload) (int64, error) {
	newLog, err := s.db.Exec(`INSERT INTO logs 
		(vehicle_id, category, title, date, description, notes, cost) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		log.VehicleId, log.Category, log.Title, log.Date, log.Description, log.Notes, log.Cost)
	if err != nil {
		return 0, err
	}

	id, err := newLog.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
