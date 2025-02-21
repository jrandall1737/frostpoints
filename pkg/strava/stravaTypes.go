package strava

import "time"

type AthleteSummary struct {
	ID                    int64         `json:"id"`
	Username              string        `json:"username"`
	ResourceState         int           `json:"resource_state"`
	Firstname             string        `json:"firstname"`
	Lastname              string        `json:"lastname"`
	City                  string        `json:"city"`
	State                 string        `json:"state"`
	Country               string        `json:"country"`
	Sex                   string        `json:"sex"`
	Premium               bool          `json:"premium"`
	CreatedAt             time.Time     `json:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at"`
	BadgeTypeID           int           `json:"badge_type_id"`
	ProfileMedium         string        `json:"profile_medium"`
	Profile               string        `json:"profile"`
	Friend                *int          `json:"friend"`   // Nullable
	Follower              *int          `json:"follower"` // Nullable
	FollowerCount         int           `json:"follower_count"`
	FriendCount           int           `json:"friend_count"`
	MutualFriendCount     int           `json:"mutual_friend_count"`
	AthleteType           int           `json:"athlete_type"`
	DatePreference        string        `json:"date_preference"`
	MeasurementPreference string        `json:"measurement_preference"`
	Clubs                 []interface{} `json:"clubs"` // Clubs is an empty array in the example, type can be adjusted if needed
	FTP                   *int          `json:"ftp"`   // Nullable
	Weight                float64       `json:"weight"`
	Bikes                 []Bike        `json:"bikes"`
	Shoes                 []Shoe        `json:"shoes"`
}

type Bike struct {
	ID            string  `json:"id"`
	Primary       bool    `json:"primary"`
	Name          string  `json:"name"`
	ResourceState int     `json:"resource_state"`
	Distance      float64 `json:"distance"`
}

type Shoe struct {
	ID            string  `json:"id"`
	Primary       bool    `json:"primary"`
	Name          string  `json:"name"`
	ResourceState int     `json:"resource_state"`
	Distance      float64 `json:"distance"`
}

type Location [2]float64

type ActivitySummary struct {
	Id                 int64          `json:"id"`
	ExternalId         string         `json:"external_id"`
	UploadId           int64          `json:"upload_id"`
	Athlete            AthleteSummary `json:"athlete"`
	Name               string         `json:"name"`
	Distance           float64        `json:"distance"`
	MovingTime         int            `json:"moving_time"`
	ElapsedTime        int            `json:"elapsed_time"`
	TotalElevationGain float64        `json:"total_elevation_gain"`
	Type               ActivityType   `json:"type"`
	Description        string         `json:"description"`

	StartDate      time.Time `json:"start_date"`
	StartDateLocal time.Time `json:"start_date_local"`

	TimeZone             string   `json:"time_zone"`
	StartLocation        Location `json:"start_latlng"`
	EndLocation          Location `json:"end_latlng"`
	City                 string   `json:"location_city"`
	State                string   `json:"location_state"`
	Country              string   `json:"location_country"`
	AchievementCount     int      `json:"achievement_count"`
	KudosCount           int      `json:"kudos_count"`
	CommentCount         int      `json:"comment_count"`
	AthleteCount         int      `json:"athlete_count"`
	PhotoCount           int      `json:"photo_count"`
	Trainer              bool     `json:"trainer"`
	Commute              bool     `json:"commute"`
	Manual               bool     `json:"manual"`
	Private              bool     `json:"private"`
	Flagged              bool     `json:"flagged"`
	GearId               string   `json:"gear_id"` // bike or pair of shoes
	AverageSpeed         float64  `json:"average_speed"`
	MaximumSpeed         float64  `json:"max_speed"`
	AverageCadence       float64  `json:"average_cadence"`
	AverageTemperature   float64  `json:"average_temp"`
	AveragePower         float64  `json:"average_watts"`
	WeightedAveragePower int      `json:"weighted_average_watts"`
	Kilojoules           float64  `json:"kilojoules"`
	DeviceWatts          bool     `json:"device_watts"`
	AverageHeartrate     float64  `json:"average_heartrate"`
	MaximumHeartrate     float64  `json:"max_heartrate"`
	Truncated            int      `json:"truncated"` // only present if activity is owned by authenticated athlete, returns 0 if not truncated by privacy zones
	HasKudoed            bool     `json:"has_kudoed"`
}

type ActivityType string

var ActivityTypes = struct {
	Ride               ActivityType
	AlpineSki          ActivityType
	BackcountrySki     ActivityType
	Hike               ActivityType
	IceSkate           ActivityType
	InlineSkate        ActivityType
	NordicSki          ActivityType
	RollerSki          ActivityType
	Run                ActivityType
	Walk               ActivityType
	Workout            ActivityType
	Snowboard          ActivityType
	Snowshoe           ActivityType
	Kitesurf           ActivityType
	Windsurf           ActivityType
	Swim               ActivityType
	VirtualRide        ActivityType
	EBikeRide          ActivityType
	WaterSport         ActivityType
	Canoeing           ActivityType
	Kayaking           ActivityType
	Rowing             ActivityType
	StandUpPaddling    ActivityType
	Surfing            ActivityType
	Crossfit           ActivityType
	Elliptical         ActivityType
	RockClimbing       ActivityType
	StairStepper       ActivityType
	WeightTraining     ActivityType
	Yoga               ActivityType
	WinterSport        ActivityType
	CrossCountrySkiing ActivityType
}{"Ride", "AlpineSki", "BackcountrySki", "Hike", "IceSkate", "InlineSkate", "NordicSki", "RollerSki",
	"Run", "Walk", "Workout", "Snowboard", "Snowshoe", "Kitesurf", "Windsurf", "Swim", "VirtualRide", "EBikeRide",

	"WaterSport", "Canoeing", "Kayaking", "Rowing", "StandUpPaddling", "Surfing",
	"Crossfit", "Elliptical", "RockClimbing", "StairStepper", "WeightTraining", "Yoga", "WinterSport", "CrossCountrySkiing",
}

type UpdatableActivity struct {
	Description string `json:"description"`
}
