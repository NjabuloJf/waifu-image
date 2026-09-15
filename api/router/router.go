package router

import (
	"log"
	"net/http"

	"github.com/NjabuloJf/waifu-image/api"
	"github.com/NjabuloJf/waifu-image/api/routes/admin"
	"github.com/NjabuloJf/waifu-image/api/routes/image"
	"github.com/NjabuloJf/waifu-image/api/routes/info"
	"github.com/NjabuloJf/waifu-image/api/routes/upload"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	echoInstance *echo.Echo
)

// init initializes the Echo instance once
func init() {
	echoInstance = echo.New()
}

// setupRouter configures all routes
func setupRouter(options api.Options) {
	e := echoInstance

	// Panic recovery middleware
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

	apiGroup := e.Group("")

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     []string{"*", options.Config.Frontend},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowCredentials: true,
	}))

	// Register routes
	image.NewRouter(options, apiGroup)
	admin.NewRouter(options, apiGroup)
	upload.NewRouter(options, apiGroup)
	info.NewRouter(options, apiGroup)

	// Endpoints route
	apiGroup.GET("/endpoints", func(c echo.Context) error {
		return c.JSON(200, options.Config.Endpoints)
	})
}

var routesInitialized = false

// HandleRequest handles individual HTTP requests for Vercel
func HandleRequest(options api.Options, w http.ResponseWriter, r *http.Request) {
	// Initialize routes only once
	if !routesInitialized {
		setupRouter(options)
		routesInitialized = true
	}

	// Serve the request through Echo
	echoInstance.ServeHTTP(w, r)
}

// New : initialize router (kept for backward compatibility, but not used on Vercel)
func New(options api.Options) {
	e := echo.New()

	// Panic recovery middleware
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

	apiGroup := e.Group("")

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     []string{"*", options.Config.Frontend},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowCredentials: true,
	}))

	// Register routes
	image.NewRouter(options, apiGroup)
	admin.NewRouter(options, apiGroup)
	upload.NewRouter(options, apiGroup)
	info.NewRouter(options, apiGroup)

	// Endpoints route
	apiGroup.GET("/endpoints", func(c echo.Context) error {
		return c.JSON(200, options.Config.Endpoints)
	})

	log.Printf("Starting server...")
	e.Logger.Fatal(e.Start(":3000"))
}
