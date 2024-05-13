package follower

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

func (s *Store) FollowUser(followerId, followingId uuid.UUID) error {
	_, err := s.db.Exec(`INSERT INTO followers (id, follower_id, following_id) VALUES (?, ?, ?)`, uuid.New(), followerId, followingId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UnfollowUser(followerId, followingId uuid.UUID) error {
	_, err := s.db.Exec(`DELETE FROM followers WHERE follower_id = ? AND following_id = ?`, followerId, followingId)
	if err != nil {
		return err
	}

	return nil
}

func scanRowIntoFollower(rows *sql.Rows) (*types.Follower, error) {
	user := new(types.Follower)

	err := rows.Scan(
		&user.ID,
		&user.FollowerID,
		&user.FollowingID,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) GetFollower(followerId, followingId uuid.UUID) (*types.Follower, error) {
	rows, err := s.db.Query(`SELECT * FROM followers WHERE follower_id = ? AND following_id = ?`, followerId, followingId)
	if err != nil {
		return nil, err
	}

	f := new(types.Follower)

	for rows.Next() {
		f, err = scanRowIntoFollower(rows)
		if err != nil {
			return nil, err
		}
	}

	if f.ID == uuid.Nil {
		return nil, nil
	}

	return f, nil
}
