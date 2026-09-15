package api

import (
	s3simple "github.com/Riku32/s3-simple"
	"github.com/NjabuloJf/waifu-image/config"
	"github.com/NjabuloJf/waifu-image/database"
)

// Options : route object
type Options struct {
	Database database.Database
	Config   config.Config
	S3       *s3simple.Session
}
