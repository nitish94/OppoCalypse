package main

import (
	"OppoCalypse/internal/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	router := gin.Default()

	// Set up sessions with secure configuration
	store := cookie.NewStore([]byte("replace-with-a-strong-secret-key-1234567890"))

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

	// Load HTML templates
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*.tmpl")

	// Setup routes
	routes.SetupRoutes(router)

	// Start the server
	port := "8081"
	fmt.Printf("Server is running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
