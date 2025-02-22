package database

type Database interface {
	AddToken(token UserToken) error
	FindTokenById(athleteId int64) *UserToken
	DeleteToken(athleteId int64) error
	Disconnect()
}
