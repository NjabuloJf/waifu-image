package main

import (
	"flag"
	"log"

	s3simple "github.com/Riku32/s3-simple"
	"github.com/NjabuloJf/waifu-image/api"
	"github.com/NjabuloJf/waifu-image/api/router"
	"github.com/NjabuloJf/waifu-image/cmd/admin"
	"github.com/NjabuloJf/waifu-image/config"
	"github.com/NjabuloJf/waifu-image/database"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	newuser := flag.Bool("newuser", false, "create an administrator")
	flag.Parse()

	// Load configuration
	conf := config.LoadConfig()
	if conf == nil {
		log.Fatalln("Failed to load config")
	}

	// Initialize database
	db := database.InitSQL(conf)
	if db == nil {
		log.Fatalln("Failed to initialize database")
	}
	defer db.Close()

	// Admin creation argument
	if *newuser {
		admin.CreateAdmin(db)
		return
	}

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
		log.Fatalf("Unable to start S3: %v", err)
	}

	options := api.Options{
		Database: db,
		Config:   conf,
		S3:       s3,
	}

	// Start the router.
	// This function should internally start the Echo server.
	router.New(options)
}
