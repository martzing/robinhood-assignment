package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Mongo      mongo
	HTTPServer httpServer
	Auth       auth
}

type mongo struct {
	URI      string `envconfig:"MONGO_URI" default:"mongodb://localhost:27017"`
	Database string `envconfig:"DB_NAME" default:"interview"`
}

type httpServer struct {
	Port int `envconfig:"PORT" default:"8080"`
}

type auth struct {
	BcryptCost int    `envconfig:"BCRYPT_COST" default:"8"`
	JwtSecret  string `envconfig:"JWT_SECRET"`
}

var cfg config

func New() {
	_ = godotenv.Load()
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("read env error : %s", err.Error())
	}
}

func Get() config {
	return cfg
}
