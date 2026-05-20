package serve

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

type ServeMode string

const (
	ModeComponent ServeMode = "component"
	ModeKit       ServeMode = "kit"
	ModeApp       ServeMode = "app"
)

type ServerConfig struct {
	Port            int
	Host            string
	Dir             string
	Mode            ServeMode
	AutoDetect      bool
	OpenBrowser     bool
	LiveReload      bool
	WebSocketPath   string
	ShutdownTimeout time.Duration
}

func DefaultConfig() *ServerConfig {
	return &ServerConfig{
		Port:            3000,
		Host:            "localhost",
		Dir:             ".",
		AutoDetect:      true,
		OpenBrowser:     true,
		LiveReload:      true,
		WebSocketPath:   "/ws",
		ShutdownTimeout: 10 * time.Second,
	}
}

type Server struct {
	config        *ServerConfig
	httpServer    *http.Server
	mux           *http.ServeMux
	watcher       *Watcher
	wsManager     *WebSocketManager
	detector      *ModeDetector
	componentMode *ComponentMode
	kitMode       *KitMode
	appMode       *AppMode
	mu            sync.RWMutex
	running       bool
}

func NewServer(config *ServerConfig) (*Server, error) {
	if config == nil {
		config = DefaultConfig()
	}

	absDir, err := filepath.Abs(config.Dir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve directory: %w", err)
	}
	config.Dir = absDir

	s := &Server{
		config:    config,
		mux:       http.NewServeMux(),
		wsManager: NewWebSocketManager(),
		detector:  NewModeDetector(absDir),
	}

	if config.AutoDetect {
		mode, err := s.detector.DetectMode()
		if err != nil {
			return nil, fmt.Errorf("failed to detect serve mode: %w", err)
		}
		config.Mode = mode
		log.Printf("Auto-detected serve mode: %s", mode)
	}

	if config.LiveReload {
		watcher, err := NewWatcher(absDir, s.handleFileChange)
		if err != nil {
			return nil, fmt.Errorf("failed to create file watcher: %w", err)
		}
		s.watcher = watcher
	}

	s.setupRoutes()

	return s, nil
}

func (s *Server) setupRoutes() {
	s.mux.HandleFunc(s.config.WebSocketPath, s.wsManager.HandleWebSocket)

	switch s.config.Mode {
	case ModeComponent:
		s.setupComponentRoutes()
	case ModeKit:
		s.setupKitRoutes()
	case ModeApp:
		s.setupAppRoutes()
	}
}

func (s *Server) setupComponentRoutes() {
	log.Println("Setting up component development routes")

	cm, err := NewComponentMode(s)
	if err != nil {
		log.Printf("Error: Failed to initialize component mode: %v", err)
		// Register error handler to prevent unhandled routes
		s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, fmt.Sprintf("Component mode failed to initialize: %v", err), http.StatusInternalServerError)
		})
		return
	}
	s.componentMode = cm

	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cm.ServeHTTP(w, r)
	})
	s.mux.HandleFunc("/preview", cm.handlePreview)
	s.mux.HandleFunc("/render", cm.handleRender)
	s.mux.HandleFunc("/reload", cm.handleReload)
}

func (s *Server) setupKitRoutes() {
	log.Println("Setting up kit development routes")

	km, err := NewKitMode(s)
	if err != nil {
		log.Printf("Error: Failed to initialize kit mode: %v", err)
		// Register error handler to prevent unhandled routes
		s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, fmt.Sprintf("Kit mode failed to initialize: %v", err), http.StatusInternalServerError)
		})
		return
	}
	s.kitMode = km

	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		km.ServeHTTP(w, r)
	})
	s.mux.HandleFunc("/test", km.handleTest)
	s.mux.HandleFunc("/helpers", km.handleHelpers)
}

func (s *Server) setupAppRoutes() {
	log.Println("Setting up app development routes")

	am, err := NewAppMode(s)
	if err != nil {
		log.Printf("Error: Failed to initialize app mode: %v", err)
		// Register error handler to prevent unhandled routes
		s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, fmt.Sprintf("App mode failed to initialize: %v", err), http.StatusInternalServerError)
		})
		return
	}
	s.appMode = am

	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		am.ServeHTTP(w, r)
	})
}

func (s *Server) handleFileChange(path string) {
	log.Printf("File changed: %s", path)

	if s.appMode != nil {
		s.appMode.HandleFileChange(path)
	}

	s.wsManager.ReloadClients(map[string]interface{}{
		"type": "reload",
		"path": path,
	})
}

func (s *Server) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("server is already running")
	}
	s.running = true
	s.mu.Unlock()

	if s.watcher != nil {
		if err := s.watcher.Start(); err != nil {
			return fmt.Errorf("failed to start file watcher: %w", err)
		}
	}

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	if !s.isPortAvailable(s.config.Port) {
		return fmt.Errorf("port %d is already in use", s.config.Port)
	}

	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errChan := make(chan error, 1)
	go func() {
		log.Printf("Starting server at http://%s (mode: %s)", addr, s.config.Mode)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	if s.config.OpenBrowser {
		go s.openBrowser(fmt.Sprintf("http://%s", addr))
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	select {
	case err := <-errChan:
		return err
	case <-sigChan:
		log.Println("\nShutting down server...")
		return s.Shutdown()
	case <-ctx.Done():
		log.Println("\nContext cancelled, shutting down server...")
		return s.Shutdown()
	}
}

func (s *Server) Shutdown() error {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return nil
	}
	s.running = false
	s.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	if s.watcher != nil {
		s.watcher.Stop()
	}

	s.wsManager.Close()

	if s.appMode != nil {
		s.appMode.Stop()
	}

	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("server shutdown failed: %w", err)
		}
	}

	log.Println("Server stopped gracefully")
	return nil
}

func (s *Server) isPortAvailable(port int) bool {
	addr := fmt.Sprintf("%s:%d", s.config.Host, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

func (s *Server) openBrowser(url string) {
	time.Sleep(500 * time.Millisecond)
	log.Printf("Please open your browser at: %s", url)
}
