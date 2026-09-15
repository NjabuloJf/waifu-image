package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Config struct {
	Port        string
	Web         Web
	DatabaseUrl string
	Storage     Storage
	Endpoints   Endpoints
	Domain      string
	Frontend    string
}

type Web struct {
	Cdn string
	Jwt string
}

type Storage struct {
	Endpoint  string
	Accesskey string
	Secretkey string
	Region    string
	Bucket    string
}

type Endpoints struct {
	Sfw  []string `json:"sfw"`
	Nsfw []string `json:"nsfw"`
}

// LoadConfig : Load config from env
func LoadConfig() Config {
	return Config{
		// PORT falls back to 3000 if Vercel doesn't provide one
		Port:        getEnvOr("PORT", "3000"),
		
		// DATABASE_URL is critical. App cannot run without it.
		DatabaseUrl: getEnv("DATABASE_URL"),
		
		Web: Web{
			Cdn: getEnvOr("CDN_URL", ""),
			// Provide a fallback JWT key so it doesn't crash, but you SHOULD set this in Vercel
			Jwt: getEnvOr("JWT_KEY", "super_secret_fallback_key_change_me"),
		},
		Storage: Storage{
			// These S3 variables are critical for your upload/image routes
			Endpoint:  getEnv("S3_ENDPOINT"),
			Accesskey: getEnv("S3_ACCESS_KEY"),
			Secretkey: getEnv("S3_SECRET_KEY"),
			Region:    getEnvOr("S3_REGION", "us-east-1"),
			Bucket:    getEnv("S3_BUCKET"),
		},
		Endpoints: Endpoints{
			// These will default to "sfw" and "nsfw" if not set
			Sfw:  strings.Split(getEnvOr("ENDPOINTS_SFW", "sfw"), ","),
			Nsfw: strings.Split(getEnvOr("ENDPOINTS_NSFW", "nsfw"), ","),
		},
		// These will fall back to safe defaults if not set
		Domain:   getEnvOr("DOMAIN", "localhost"),
		Frontend: getEnvOr("FRONTEND_URL", "*"),
	}
}

// getEnv : Returns error if a critical variable is missing
func getEnv(key string) string {
	value, set := os.LookupEnv(key)
	if !set || value == "" {
		log.Fatalln(fmt.Sprintf("Required config variable %s was missing", key))
	}
	return value
}

// getEnvOr : Returns a fallback value if the variable is missing
func getEnvOr(key string, fallback string) string {
	value, set := os.LookupEnv(key)
	if !set || value == "" {
		return fallback
	}
	return value
}
