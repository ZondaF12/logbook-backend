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
}

type FollowerStore interface {
	FollowUser(followerId, followingId int) error
	UnfollowUser(followerId, followingId int) error
	GetFollower(followerId, followingId int) (*Follower, error)
}

type VehicleStore interface {
	GetVehicleByID(id int) (*Vehicle, error)
	GetAuthenticatedUserVehicles(userID int) ([]Vehicle, error)
	GetVehicleByRegistration(registration string) (*Vehicle, error)
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

type Vehicle struct {
	ID            string   `json:"id"`
	UserID        string   `json:"user_id"`
	Registration  string   `json:"registration"`
	Color         string   `json:"color"`
	Description   string   `json:"description"`
	EngineSize    uint16   `json:"engine_size"`
	Images        []string `json:"images"`
	InsuranceDate string   `json:"insurance_date"`
	Make          string   `json:"make"`
	Model         string   `json:"model"`
	MotDate       string   `json:"mot_date"`
	Nickname      string   `json:"nickname"`
	Registered    string   `json:"registered"`
	ServiceDate   string   `json:"service_date"`
	TaxDate       string   `json:"tax_date"`
	Year          uint16   `json:"year"`
}

type NewVehicle struct {
	Registration  string   `json:"registration"`
	Color         string   `json:"color"`
	Description   string   `json:"description"`
	EngineSize    uint16   `json:"engine_size"`
	Images        []string `json:"images"`
	InsuranceDate string   `json:"insurance_date"`
	Make          string   `json:"make"`
	Model         string   `json:"model"`
	MotDate       string   `json:"mot_date"`
	Nickname      string   `json:"nickname"`
	Registered    string   `json:"registered"`
	ServiceDate   string   `json:"service_date"`
	TaxDate       string   `json:"tax_date"`
	Year          uint16   `json:"year"`
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

type NewVehiclePostData struct {
	Registration string   `json:"registration"`
	Images       []string `json:"images"`
	Nickname     string   `json:"nickname"`
	Model        string   `json:"model"`
	Description  string   `json:"description"`
}
