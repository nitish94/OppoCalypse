package main

import (
	"OppoCalypse/internal/middleware"
	"OppoCalypse/internal/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}

	// Get port from env, default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Get session secret from env
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "default-secret-key-change-in-production"
	}

	// Initialize Gin router
	router := gin.Default()

	// Set up sessions with secure configuration
	store := cookie.NewStore([]byte(sessionSecret))

	// Set session options
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	// Initialize session middleware with cookie settings
	sessionMiddleware := sessions.Sessions("oppocalypse_session", store)

	// Apply the session middleware
	router.Use(sessionMiddleware)

	// Apply logging middleware
	router.Use(middleware.Logger())

	// Load HTML templates
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*.tmpl")

	// Setup routes
	routes.SetupRoutes(router)

	// Start the server
	fmt.Printf("Server is running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
