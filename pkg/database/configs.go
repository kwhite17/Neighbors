package database

type DbConfig struct {
	Driver          string
	Host            string
	DevelopmentMode bool
}

var SQLITE3 = &DbConfig{
	Driver:          "sqlite3",
	Host:            "file::memory:?mode=memory&cache=shared",
	DevelopmentMode: true,
}

func BuildConfig(driver string, host string, developmentMode bool) *DbConfig {
	return &DbConfig{Driver: driver, Host: host, DevelopmentMode: developmentMode}
}
