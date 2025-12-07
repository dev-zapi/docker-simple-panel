package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/dev-zapi/docker-simple-panel/config"
	"github.com/dev-zapi/docker-simple-panel/database"
	"github.com/dev-zapi/docker-simple-panel/docker"
	"github.com/dev-zapi/docker-simple-panel/handlers"
	"github.com/dev-zapi/docker-simple-panel/middleware"
)

func main() {
	// Load configuration from environment
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.NewDB(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connected successfully")

	// Helper function to load config from database or use default
	loadConfigString := func(key, defaultValue string) string {
		if value, err := db.GetConfig(key); err == nil && value != "" {
			log.Printf("Loaded %s from database: %s", key, value)
			return value
		}
		return defaultValue
	}

	loadConfigBool := func(key string, defaultValue bool) bool {
		if value, err := db.GetConfig(key); err == nil && value != "" {
			result := value == "true"
			log.Printf("Loaded %s from database: %v", key, result)
			return result
		}
		return defaultValue
	}

	// Load persisted configs from database or use environment defaults
	dockerSocket := loadConfigString("docker_socket", cfg.DockerSocket)
	disableRegistration := loadConfigBool("disable_registration", cfg.DisableRegistration)
	logLevel := config.ParseLogLevel(loadConfigString("log_level", cfg.LogLevel.String()))

	// Initialize configuration manager
	configManager := config.NewManager(dockerSocket, disableRegistration, logLevel)
	log.Printf("Log level set to: %s", logLevel.String())

	// Initialize Docker manager
	dockerManager, err := docker.NewManager(dockerSocket)
	if err != nil {
		log.Fatalf("Failed to create Docker manager: %v", err)
	}
	defer dockerManager.Close()
	log.Println("Docker manager initialized successfully")

	// Test Docker connection
	ctx := context.Background()
	if err := dockerManager.Ping(ctx); err != nil {
		log.Printf("Warning: Docker daemon not accessible: %v", err)
	} else {
		log.Println("Docker daemon is accessible")
	}

	// Set Docker socket change callback
	configManager.SetDockerSocketChangeCallback(func(newSocket string) error {
		return dockerManager.RestartWithSocket(newSocket)
	})

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret, configManager)
	dockerHandler := handlers.NewDockerHandler(dockerManager)
	configHandler := handlers.NewConfigHandler(configManager, db)
	userHandler := handlers.NewUserHandler(db)

	// Setup router
	router := mux.NewRouter()

	// Apply CORS middleware
	router.Use(middleware.CORS)

	// Apply logging middleware
	router.Use(middleware.Logging(configManager))

	// Public routes
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true,"message":"Server is running"}`))
	}).Methods("GET")

	// Serve OpenAPI specification
	router.HandleFunc("/api/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		file, err := os.Open("./openapi.json")
		if err != nil {
			log.Printf("Failed to open OpenAPI specification: %v", err)
			http.Error(w, "OpenAPI specification not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		if _, err := io.Copy(w, file); err != nil {
			log.Printf("Failed to serve OpenAPI specification: %v", err)
		}
	}).Methods("GET")

	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")

	// Public config endpoint (no auth required) - returns only registration status
	router.HandleFunc("/api/config/public", configHandler.GetPublicConfig).Methods("GET")

	// Protected routes - require JWT authentication
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.JWTAuth(cfg.JWTSecret))

	// Docker container routes
	protected.HandleFunc("/containers", dockerHandler.ListContainers).Methods("GET")
	protected.HandleFunc("/containers/{id}", dockerHandler.GetContainer).Methods("GET")
	protected.HandleFunc("/containers/{id}/start", dockerHandler.StartContainer).Methods("POST")
	protected.HandleFunc("/containers/{id}/stop", dockerHandler.StopContainer).Methods("POST")
	protected.HandleFunc("/containers/{id}/restart", dockerHandler.RestartContainer).Methods("POST")
	protected.HandleFunc("/containers/{id}/logs/stream", dockerHandler.StreamContainerLogs).Methods("GET")
	protected.HandleFunc("/docker/health", dockerHandler.HealthCheck).Methods("GET")

	// Docker volume routes
	protected.HandleFunc("/volumes", dockerHandler.ListVolumes).Methods("GET")
	protected.HandleFunc("/volumes/{name}", dockerHandler.DeleteVolume).Methods("DELETE")

	// System configuration routes
	protected.HandleFunc("/config", configHandler.GetConfig).Methods("GET")
	protected.HandleFunc("/config", configHandler.UpdateConfig).Methods("PUT", "PATCH")

	// User management routes
	protected.HandleFunc("/users", userHandler.ListUsers).Methods("GET")
	protected.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	protected.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Static file serving (if configured)
	if cfg.StaticPath != "" {
		// Check if the static path exists
		if _, err := os.Stat(cfg.StaticPath); os.IsNotExist(err) {
			log.Printf("Warning: Static path does not exist: %s", cfg.StaticPath)
		} else {
			log.Printf("Serving static files from: %s", cfg.StaticPath)
			fs := http.FileServer(http.Dir(cfg.StaticPath))
			router.PathPrefix("/").Handler(fs)
		}
	}

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
