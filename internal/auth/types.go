package auth

import "time"

type Athlete struct {
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
