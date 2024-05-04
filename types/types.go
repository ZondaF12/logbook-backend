package types

import "time"

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
}

type ProfileStore interface {
	GetProfileByUserId(userId int) (*Profile, error)
	CreateProfile(Profile) error
	UpdateAvatar(userId int, avatar string) error
}

type FollowerStore interface {
	FollowUser(followerId, followingId int) error
	UnfollowUser(followerId, followingId int) error
	GetFollower(followerId, followingId int) (*Follower, error)
}

type GarageStore interface {
	GetVehicleByID(id int) (*Vehicle, error)
	GetAuthenticatedUserVehicles(userID int) ([]*Vehicle, error)
	GetVehicleByRegistration(userId int, registration string) (*Vehicle, error)
	AddUserVehicle(userID int, vehicle NewVehiclePostData) error
	CheckVehicleAdded(userId int, registration string) (bool, error)
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
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
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

type Profile struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	Avatar    string `json:"avatar"`
	Public    bool   `json:"public"`
	Followers int    `json:"followers"`
	Following int    `json:"following"`
}

type CreateProfilePayload struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Name     string `json:"name" validate:"required,min=3,max=100"`
}

type Follower struct {
	ID          int       `json:"id"`
	FollowerID  int       `json:"follower_i"`
	FollowingID int       `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type FollowUserPayload struct {
	UserID int `json:"user_id" validate:"required"`
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
	ID            string    `json:"id,omitempty"`
	UserID        string    `json:"user_id,omitempty"`
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
