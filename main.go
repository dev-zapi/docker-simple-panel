package main

import (
	"context"
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

	// Load persisted configs from database or use environment defaults
	dockerSocket := cfg.DockerSocket
	if persistedSocket, err := db.GetConfig("docker_socket"); err == nil && persistedSocket != "" {
		dockerSocket = persistedSocket
		log.Printf("Loaded Docker socket from database: %s", dockerSocket)
	}

	disableRegistration := cfg.DisableRegistration
	if persistedReg, err := db.GetConfig("disable_registration"); err == nil && persistedReg != "" {
		disableRegistration = persistedReg == "true"
		log.Printf("Loaded registration setting from database: %v", disableRegistration)
	}

	// Initialize configuration manager
	configManager := config.NewManager(dockerSocket, disableRegistration)

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

	// Setup router
	router := mux.NewRouter()

	// Apply CORS middleware
	router.Use(middleware.CORS)

	// Public routes
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true,"message":"Server is running"}`))
	}).Methods("GET")
	
	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")

	// Protected routes - require JWT authentication
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.JWTAuth(cfg.JWTSecret))

	// Docker container routes
	protected.HandleFunc("/containers", dockerHandler.ListContainers).Methods("GET")
	protected.HandleFunc("/containers/{id}", dockerHandler.GetContainer).Methods("GET")
	protected.HandleFunc("/containers/{id}/start", dockerHandler.StartContainer).Methods("POST")
	protected.HandleFunc("/containers/{id}/stop", dockerHandler.StopContainer).Methods("POST")
	protected.HandleFunc("/containers/{id}/restart", dockerHandler.RestartContainer).Methods("POST")
	protected.HandleFunc("/docker/health", dockerHandler.HealthCheck).Methods("GET")

	// System configuration routes
	protected.HandleFunc("/config", configHandler.GetConfig).Methods("GET")
	protected.HandleFunc("/config", configHandler.UpdateConfig).Methods("PUT", "PATCH")

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
