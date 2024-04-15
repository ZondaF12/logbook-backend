package types

import "time"

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByID(id int) (*User, error)
	UpdateUser(userId int, payload UpdateUserPayload) error
	CreateUser(User) error
}

type RegisterUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=100"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserPayload struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Bio      string `json:"bio,omitempty" validate:"omitempty"`
	Public   bool   `json:"public"`
}

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Bio       string    `json:"bio"`
	Public    bool      `json:"public"`
	CreatedAt time.Time `json:"created_at"`
}
