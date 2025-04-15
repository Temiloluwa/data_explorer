package app

// This file could contain functions to set up and run the Server application (gRPC + REST),
// encapsulating the initialization logic currently in cmd/server/main.go.

/* Example Structure:
import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	pb "MODULE_PATH_PLACEHOLDER/api/proto/dataexplorer/v1"
	"MODULE_PATH_PLACEHOLDER/internal/config"
	dsc "MODULE_PATH_PLACEHOLDER/internal/datasources/common"
	"MODULE_PATH_PLACEHOLDER/internal/datasources/local"
	"MODULE_PATH_PLACEHOLDER/internal/datasources/web"
	"MODULE_PATH_PLACEHOLDER/internal/platform/logger"
	"MODULE_PATH_PLACEHOLDER/internal/qengine"
	grpchandler "MODULE_PATH_PLACEHOLDER/internal/transport/grpc"
	restgateway "MODULE_PATH_PLACEHOLDER/internal/transport/rest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServerApp struct {
	Log    logger.Logger
	Cfg    *config.Config
	Engine qengine.Engine
	// Add other shared components like DB connections if needed
}

func NewServerApp(log logger.Logger, cfg *config.Config) (*ServerApp, error) {
	// Initialize Data Sources (similar to NewCLIApp)
	var availableSources []dsc.DataSource
	if cfg.DataSources.WebSearch.Enable {
		webSource, _ := web.NewWebSearchSource(log, cfg.DataSources.WebSearch.APIKey)
        if webSource != nil { availableSources = append(availableSources, webSource) }
	}
	if cfg.DataSources.LocalFiles.Enable {
		localSource, _ := local.NewLocalFileSource(log, cfg.DataSources.LocalFiles.WatchPaths)
        if localSource != nil { availableSources = append(availableSources, localSource) }
	}
	// ... initialize other sources

	if len(availableSources) == 0 {
		log.Warnf("Server App: No data sources enabled or initialized!")
	}

	// Initialize Core Engine
	engine, err := qengine.NewEngine(log, cfg, availableSources)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize QA engine: %w", err)
	}

	return &ServerApp{
		Log:    log,
		Cfg:    cfg,
		Engine: engine,
	}, nil
}

// Run starts the gRPC and REST gateway servers and handles graceful shutdown.
func (a *ServerApp) Run(ctx context.Context) error {
	a.Log.Infof("Starting Data Explorer Server App...")
	ctx, cancel := context.WithCancel(ctx) // Create cancellable context for shutdown
	defer cancel()

	var wg sync.WaitGroup // WaitGroup to manage server goroutines

	// --- Setup gRPC Server ---
	grpcServer := grpc.NewServer( /* Add interceptors here * / )
	grpcAPIHandler := grpchandler.NewDataExplorerServer(a.Engine, a.Log)
	pb.RegisterDataExplorerServiceServer(grpcServer, grpcAPIHandler)
	reflection.Register(grpcServer) // Enable gRPC reflection

	grpcAddr := a.Cfg.Server.GRPCPort
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		a.Log.Errorf("Failed to listen on gRPC address %s: %v", grpcAddr, err)
		return fmt.Errorf("gRPC listen failed: %w", err)
	}

	// Start gRPC server in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.Log.Infof("gRPC Server listening on %s", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			// Log error unless it's ErrServerClosed during shutdown
			a.Log.Errorf("gRPC Server failed: %v", err)
			cancel() // Trigger shutdown of other components if gRPC fails critically
		}
		a.Log.Infof("gRPC Server stopped.")
	}()

	// --- Setup REST Gateway ---
	restAddr := a.Cfg.Server.RESTPort
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Pass the cancellable context to the gateway runner
		if err := restgateway.RunRESTGateway(ctx, a.Log, grpcAddr, restAddr); err != nil {
			a.Log.Errorf("REST Gateway failed: %v", err)
			// Optional: Trigger shutdown if REST gateway fails critically?
			// cancel()
		}
		a.Log.Infof("REST Gateway stopped.")
	}()


	// --- Graceful Shutdown Handling ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		a.Log.Infof("Received shutdown signal: %v. Initiating graceful shutdown...", sig)
	case <-ctx.Done(): // Shutdown triggered by context cancellation (e.g., internal error)
		a.Log.Infof("Shutdown initiated by context cancellation.")
	}

    // Trigger context cancellation for components listening to it (like REST gateway)
    cancel()

	// Graceful stop for gRPC server
	// Give it a deadline (adjust as needed)
	shutdownTimeout := 5 * time.Second
	stopped := make(chan struct{})
	go func() {
		a.Log.Infof("Attempting graceful shutdown of gRPC server (timeout: %s)...", shutdownTimeout)
		grpcServer.GracefulStop()
		close(stopped)
	}()

	// Wait for graceful stop or timeout
	select {
	case <-time.After(shutdownTimeout):
		a.Log.Warnf("gRPC server graceful shutdown timed out after %s. Forcing stop.", shutdownTimeout)
		grpcServer.Stop() // Force stop if graceful stop fails
	case <-stopped:
		a.Log.Infof("gRPC server gracefully stopped.")
	}

	// Wait for all server goroutines to finish (including REST gateway which should respond to ctx cancellation)
	a.Log.Infof("Waiting for server goroutines to complete...")
	wg.Wait()
	a.Log.Infof("Server shutdown complete.")
	return nil // Indicate successful shutdown
}


// RunServer is the entry point called by cmd/server/main.go
func RunServer() {
	// Basic initial logger
	tempLog, _ := logger.New(logger.Config{Level: "info"})

	cfg, err := config.LoadConfig(tempLog, "") // Load config early
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	log, err := logger.New(cfg.Log)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}

	app, err := NewServerApp(log, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Server application: %v", err)
	}

	// Run the application with a background context
	if err := app.Run(context.Background()); err != nil {
		log.Errorf("Application run failed: %v", err)
		os.Exit(1) // Exit with error if Run returns an error
	}
}

*/
