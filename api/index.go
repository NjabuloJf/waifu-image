package main

import (
	"net/http"

	s3simple "github.com/Riku32/s3-simple"
	"github.com/NjabuloJf/waifu-image/api"
	"github.com/NjabuloJf/waifu-image/api/router"
	"github.com/NjabuloJf/waifu-image/config"
	"github.com/NjabuloJf/waifu-image/database"
	_ "github.com/joho/godotenv/autoload"
)

var (
	options api.Options
)

// init runs once when the function is first invoked
func init() {
	// Load configuration
	conf := config.LoadConfig()

	// Initialize database
	db := database.InitSQL(conf)

	// Initiate S3
	s3, err := s3simple.New(s3simple.Config{
		Region:   conf.Storage.Region,
		Endpoint: conf.Storage.Endpoint,
		Bucket:   conf.Storage.Bucket,
		Credentials: s3simple.Credentials{
			Accesskey: conf.Storage.Accesskey,
			Secretkey: conf.Storage.Secretkey,
		},
	})

	if err != nil {
		panic("Unable to start S3: " + err.Error())
	}

	options = api.Options{
		Database: db,
		Config:   conf,
		S3:       s3,
	}
}

// Handler is the Vercel serverless function handler
func Handler(w http.ResponseWriter, r *http.Request) {
	router.HandleRequest(options, w, r)
}
