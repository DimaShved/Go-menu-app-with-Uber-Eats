package config

type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

type AppConfig struct {
	PORT string `envconfig:"PORT" default:"8080"`
}

type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	User     string `envconfig:"DB_USER" required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	DbName   string `envconfig:"DB_NAME" required:"true"`
	Port     string `envconfig:"DB_PORT" default:"5432"`
}
