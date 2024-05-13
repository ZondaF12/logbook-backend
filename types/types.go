package types

import (
	"time"

	"github.com/google/uuid"
)

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id uuid.UUID) (*User, error)
	CreateUser(User) error
}

type ProfileStore interface {
	GetProfileByUserId(userId uuid.UUID) (*Profile, error)
	CreateProfile(Profile) error
	UpdateAvatar(userId uuid.UUID, avatar string) error
}

type FollowerStore interface {
	FollowUser(followerId, followingId uuid.UUID) error
	UnfollowUser(followerId, followingId uuid.UUID) error
	GetFollower(followerId, followingId uuid.UUID) (*Follower, error)
}

type GarageStore interface {
	GetVehicleByID(id uuid.UUID) (*Vehicle, error)
	GetAuthenticatedUserVehicles(userID uuid.UUID) ([]*Vehicle, error)
	GetVehicleByRegistration(userId uuid.UUID, registration string) (*Vehicle, error)
	AddUserVehicle(userID uuid.UUID, vehicle NewVehiclePostData) (uuid.UUID, error)
	CheckVehicleAdded(userId uuid.UUID, registration string) (bool, error)
}

type MediaStore interface {
	AddNewVehicleMedia(Media) error
	AddNewLogMedia(Media) error
}

type LogbookStore interface {
	CreateLog(CreateLogPayload) (uuid.UUID, error)
	GetLogsByVehicleId(vehicleId uuid.UUID) ([]*Log, error)
}

type RegisterAuthPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=100"`
}

type LoginAuthPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Auth struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Bio       string    `json:"bio"`
	Public    bool      `json:"public"`
	CreatedAt time.Time `json:"created_at"`
}

type Profile struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	Avatar    string    `json:"avatar"`
	Public    bool      `json:"public"`
	Followers int       `json:"followers"`
	Following int       `json:"following"`
}

type CreateProfilePayload struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Name     string `json:"name" validate:"required,min=3,max=100"`
}

type Follower struct {
	ID          uuid.UUID `json:"id"`
	FollowerID  uuid.UUID `json:"follower_i"`
	FollowingID uuid.UUID `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type FollowUserPayload struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

type NewVehiclePostData struct {
	Registration string   `json:"registration"`
	Model        string   `json:"model"`
	Description  string   `json:"description"`
	Images       []string `json:"images"`
	Make         string   `json:"make"`
	Year         uint16   `json:"year"`
	EngineSize   uint16   `json:"engine_size"`
	Color        string   `json:"color"`
	Registered   string   `json:"registered"`
	TaxDate      string   `json:"tax_date"`
	MotDate      string   `json:"mot_date"`
	Nickname     string   `json:"nickname"`
}

type Vehicle struct {
	ID            uuid.UUID `json:"id,omitempty"`
	UserID        uuid.UUID `json:"user_id,omitempty"`
	Registration  string    `json:"registration,omitempty"`
	Color         string    `json:"color,omitempty"`
	Description   string    `json:"description,omitempty"`
	EngineSize    uint16    `json:"engine_size,omitempty"`
	Make          string    `json:"make,omitempty"`
	Model         string    `json:"model,omitempty"`
	MotDate       string    `json:"mot_date,omitempty"`
	Registered    string    `json:"registered,omitempty"`
	InsuranceDate string    `json:"insurance_date,omitempty"`
	ServiceDate   string    `json:"service_date,omitempty"`
	TaxDate       string    `json:"tax_date,omitempty"`
	Year          uint16    `json:"year,omitempty"`
	Mileage       uint32    `json:"mileage,omitempty"`
	Nickname      string    `json:"nickname,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	Images        string    `json:"images,omitempty"`
}

type VehicleInfoRequestData struct {
	Registration string `json:"registration"`
	Color        string `json:"color"`
	EngineSize   uint16 `json:"engine_size"`
	Make         string `json:"make"`
	Model        string `json:"model"`
	TaxDate      string `json:"tax_date"`
	MotDate      string `json:"mot_date"`
	Registered   string `json:"registered"`
	Year         uint16 `json:"year"`
}

type VehicleData struct {
	RegistrationNumber       string `json:"registrationNumber"`
	TaxStatus                string `json:"taxStatus"`
	TaxDueDate               string `json:"taxDueDate"`
	ArtEndDate               string `json:"artEndDate"`
	MotStatus                string `json:"motStatus"`
	Make                     string `json:"make"`
	YearOfManufacture        int    `json:"yearOfManufacture"`
	EngineCapacity           int    `json:"engineCapacity"`
	Co2Emissions             int    `json:"co2Emissions"`
	FuelType                 string `json:"fuelType"`
	MarkedForExport          bool   `json:"markedForExport"`
	Colour                   string `json:"colour"`
	TypeApproval             string `json:"typeApproval"`
	RevenueWeight            int    `json:"revenueWeight"`
	EuroStatus               string `json:"euroStatus"`
	DateOfLastV5CIssued      string `json:"dateOfLastV5CIssued"`
	MotExpiryDate            string `json:"motExpiryDate"`
	Wheelplan                string `json:"wheelplan"`
	MonthOfFirstRegistration string `json:"monthOfFirstRegistration"`
}

type MotData []struct {
	Registration      string     `json:"registration"`
	Make              string     `json:"make"`
	Model             string     `json:"model"`
	FirstUsedDate     string     `json:"firstUsedDate"`
	FuelType          string     `json:"fuelType"`
	PrimaryColour     string     `json:"primaryColour"`
	MotTestExpiryDate string     `json:"MotTestExpiryDate"`
	MotTests          []MotTests `json:"motTests"`
}

type MotTests struct {
	CompletedDate  string `json:"completedDate"`
	TestResult     string `json:"testResult"`
	ExpiryDate     string `json:"expiryDate"`
	OdometerValue  string `json:"odometerValue"`
	OdometerUnit   string `json:"odometerUnit"`
	MotTestNumber  string `json:"motTestNumber"`
	RfrAndComments []any  `json:"rfrAndComments"`
}

type Media struct {
	ID         *uuid.UUID `json:"id"`
	Filename   *string    `json:"filename"`
	FileType   *string    `json:"file_type"`
	S3Location *string    `json:"s3_location"`
	UploadedAt *time.Time `json:"uploaded_at"`
	UserID     *uuid.UUID `json:"user_id,omitempty"`
	VehicleID  *uuid.UUID `json:"vehicle_id,omitempty"`
	LogID      *uuid.UUID `json:"log_id,omitempty"`
}

type CreateLogPayload struct {
	VehicleId   uuid.UUID `json:"vehicle_id"`
	Title       string    `json:"title" validate:"required,min=3,max=100"`
	Category    int       `json:"category" validate:"required"`
	Date        string    `json:"date" validate:"required"`
	Description string    `json:"description"`
	Notes       string    `json:"notes"`
	Cost        float32   `json:"cost"`
}

type Log struct {
	ID          uuid.UUID   `json:"id"`
	VehicleID   uuid.UUID   `json:"vehicle_id"`
	Title       string      `json:"title"`
	Category    int         `json:"category"`
	Date        string      `json:"date"`
	Description string      `json:"description"`
	Notes       string      `json:"notes"`
	Cost        float32     `json:"cost"`
	CreatedAt   time.Time   `json:"created_at"`
	Media       []*LogMedia `json:"media"`
}

type LogMedia struct {
	Filename   *string `json:"filename"`
	FileType   *string `json:"file_type"`
	S3Location *string `json:"s3_location"`
}
