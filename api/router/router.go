package router

import (
	"log"
	"net/http"
	"os"

	"github.com/NjabuloJf/waifu-image/api"
	"github.com/NjabuloJf/waifu-image/api/routes/admin"
	"github.com/NjabuloJf/waifu-image/api/routes/image"
	"github.com/NjabuloJf/waifu-image/api/routes/info"
	"github.com/NjabuloJf/waifu-image/api/routes/upload"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// New : initialize router
func New(options api.Options) {
	e := echo.New()

	// 👇 ADD PANIC RECOVERY MIDDLEWARE HERE 👇
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("PANIC: %v", r)
					c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
				}
			}()
			return next(c)
		}
	})
	// 👆 END OF PANIC RECOVERY MIDDLEWARE 👆

	api := e.Group("") // Root URL for the API location
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     []string{"*", options.Config.Frontend},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowCredentials: true,
	}))

	image.NewRouter(options, api)
	admin.NewRouter(options, api)
	upload.NewRouter(options, api)
	info.NewRouter(options, api)

	api.GET("/endpoints", func(c echo.Context) error {
		return c.JSON(200, options.Config.Endpoints)
	})

	// 👇 VERIFY VERCEL PORT 👇
	// Vercel sets a PORT environment variable. We must use it.
	port := os.Getenv("PORT")
	if port == "" {
		port = options.Config.Port // Fallback to your config if not on Vercel
	}

	log.Printf("Starting server on port %s...", port)
	e.Logger.Fatal(e.Start(":" + port))
}
