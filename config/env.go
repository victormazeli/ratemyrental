package config

import (
	"github.com/spf13/viper"
	"log"
)

type Env struct {
	ServerPort       string `mapstructure:"SERVER_PORT"`
	DBHost           string `mapstructure:"DB_HOST"`
	DBPort           string `mapstructure:"DB_PORT"`
	DBUser           string `mapstructure:"DB_USER"`
	DBPass           string `mapstructure:"DB_PASS"`
	DBName           string `mapstructure:"DB_NAME"`
	JwtKey           string `mapstructure:"JWT_KEY"`
	RedisUrl         string `mapstructure:"REDIS_URL"`
	BaseUrl          string `mapstructure:"BASE_URL"`
	MailGunDomain    string `mapstructure:"MAILGUN_DOMAIN"`
	MailGunApiKey    string `mapstructure:"MAILGUN_API_KEY"`
	CloudinaryName   string `mapstructure:"CLOUDINARY_NAME"`
	CloudinaryApiKey string `mapstructure:"CLOUDINARY_APIKEY"`
	CloudinarySecret string `mapstructure:"CLOUDINARY_SECRET"`
}

func NewEnv(envType string) *Env {
	env := Env{}
	//viper.SetConfigFile("development.env")
	viper.SetConfigType("env")
	viper.SetConfigName(envType)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Can't find the file %s.env : %v", envType, err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	return &env
}
