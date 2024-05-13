package profile

import (
	"database/sql"
	"fmt"

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

func (s *Store) GetProfileByUserId(userId uuid.UUID) (*types.Profile, error) {
	rows, err := s.db.Query(`
	SELECT
		p.*,
		(SELECT COUNT(*) FROM followers f WHERE f.following_id = p.user_id) AS followers,
		(SELECT COUNT(*) FROM followers f WHERE f.follower_id = p.user_id) AS following
	FROM
		profiles p
	WHERE 
		p.user_id = ?`, userId)
	if err != nil {
		return nil, err
	}

	u := new(types.Profile)

	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == uuid.Nil {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func scanRowIntoUser(rows *sql.Rows) (*types.Profile, error) {
	user := new(types.Profile)

	err := rows.Scan(
		&user.ID,
		&user.UserID,
		&user.Username,
		&user.Name,
		&user.Bio,
		&user.Avatar,
		&user.Public,
		&user.Followers,
		&user.Following,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) CreateProfile(u types.Profile) error {
	_, err := s.db.Exec("INSERT INTO profiles (id, user_id, username, name) VALUES (?, ?, ?, ?)", uuid.New(), u.UserID, u.Username, u.Name)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateAvatar(userId uuid.UUID, avatar string) error {
	_, err := s.db.Exec("UPDATE profiles SET avatar = ? WHERE user_id = ?", avatar, userId)
	if err != nil {
		return err
	}

	return nil
}
