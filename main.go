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
	"github.com/dev-zapi/docker-simple-panel/docker"
	"github.com/dev-zapi/docker-simple-panel/handlers"
	"github.com/dev-zapi/docker-simple-panel/middleware"
)

func main() {
	// Load configuration from YAML file
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Configuration loaded successfully")

	// Initialize configuration manager
	configManager := config.NewManager(cfg)
	log.Printf("Log level set to: %s", cfg.Logging.Level)
	log.Printf("Session max timeout set to: %d hours", cfg.Server.SessionMaxTimeout)
	log.Printf("Username configured: %s", cfg.Username)

	// Initialize Docker manager
	dockerManager, err := docker.NewManager(cfg.Docker.Socket)
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
	authHandler := handlers.NewAuthHandler(configManager, cfg.Server.JWTSecret)
	dockerHandler := handlers.NewDockerHandler(dockerManager, configManager)
	configHandler := handlers.NewConfigHandler(configManager)

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

	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")

	// Protected routes - require JWT authentication
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.JWTAuth(cfg.Server.JWTSecret))

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
	protected.HandleFunc("/volumes/{name}/files", dockerHandler.ExploreVolumeFiles).Methods("GET")
	protected.HandleFunc("/volumes/{name}/file", dockerHandler.ReadVolumeFile).Methods("GET")
	protected.HandleFunc("/volumes/{name}", dockerHandler.DeleteVolume).Methods("DELETE")

	// System configuration routes
	protected.HandleFunc("/config", configHandler.GetConfig).Methods("GET")
	protected.HandleFunc("/config", configHandler.UpdateConfig).Methods("PUT", "PATCH")

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
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
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
