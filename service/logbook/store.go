package logbook

import (
	"database/sql"

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

func (s *Store) CreateLog(log types.CreateLogPayload) (uuid.UUID, error) {
	newLogId := uuid.New()
	_, err := s.db.Exec(`INSERT INTO logs 
		(id, vehicle_id, category, title, date, description, notes, cost) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		newLogId, log.VehicleId, log.Category, log.Title, log.Date, log.Description, log.Notes, log.Cost)
	if err != nil {
		return uuid.Nil, err
	}

	return newLogId, nil
}

func scanRowIntoLogWithMedia(rows *sql.Rows) (*types.Log, *types.LogMedia, error) {
	var logbook types.Log
	var media types.LogMedia

	err := rows.Scan(
		&logbook.ID,
		&logbook.VehicleID,
		&logbook.Category,
		&logbook.Title,
		&logbook.Date,
		&logbook.Description,
		&logbook.Notes,
		&logbook.Cost,
		&logbook.CreatedAt,
		&media.Filename,
		&media.FileType,
		&media.S3Location,
	)
	if err != nil {
		return nil, nil, err
	}

	return &logbook, &media, nil
}

func (s *Store) GetLogsByVehicleId(vehicleId uuid.UUID) ([]*types.Log, error) {
	rows, err := s.db.Query(`
		SELECT
			logs.*,
			media.filename,
			media.file_type,
			media.s3_location
		FROM logs
		LEFT JOIN media
			ON logs.id = media.log_id
		WHERE logs.vehicle_id = ?
		ORDER BY logs.created_at DESC`, vehicleId)
	if err != nil {
		return nil, err
	}

	logs := make(map[string]*types.Log)

	for rows.Next() {
		log, media, err := scanRowIntoLogWithMedia(rows)
		if err != nil {
			return nil, err
		}

		logID := log.ID.String()

		// If log does not exist in map, create it
		if _, ok := logs[logID]; !ok {
			log.Media = []*types.LogMedia{}
			logs[logID] = log
		}

		// If media ID is not nil, add media to log
		if media.Filename != nil {
			logs[logID].Media = append(logs[logID].Media, media)
		}
	}

	// Convert map to slice
	logSlice := make([]*types.Log, 0, len(logs))
	for _, log := range logs {
		logSlice = append(logSlice, log)
	}

	return logSlice, nil
}
