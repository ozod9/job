package models

type Config struct {
	Application application
	Database    database
}

type database struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     int
}

type application struct {
	Name    string
	Version string
	Host    string
}
