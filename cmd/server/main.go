package main

import (
	"context" // Keep context import
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync" // Added for WaitGroup
	"syscall"
	"time"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/yourusername/data-explorer/api/proto/dataexplorer/v1" // Adjust
	"github.com/yourusername/data-explorer/internal/config"       // Adjust
	dsc "github.com/yourusername/data-explorer/internal/datasources/common" // Adjust
	// Import specific datasource packages
	"github.com/yourusername/data-explorer/internal/datasources/local" // Adjust
	"github.com/yourusername/data-explorer/internal/datasources/web"  // Adjust
	"github.com/yourusername/data-explorer/internal/platform/logger"    // Adjust
	"github.com/yourusername/data-explorer/internal/qengine"      // Adjust
	grpchandler "github.com/yourusername/data-explorer/internal/transport/grpc" // Adjust
	restgateway "github.com/yourusername/data-explorer/internal/transport/rest" // Adjust for refactored gateway
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // Use secure creds in production!
	"google.golang.org/grpc/reflection"           // For gRPC reflection
)

func main() {
	// --- Basic Logger First ---
	// Use a temporary basic logger until config is loaded
	tempLog, _ := logger.New(logger.Config{Level: "info", Format: "console"}) // Ensure console format for startup
	tempLog.Infof("Starting Data Explorer Server...")

	// --- Configuration ---
	// Consider using flags package for command-line arguments like --config
	var configPath string // Example: flag.StringVar(&configPath, "config", "", "Path to config file")
	// flag.Parse()

	cfg, err := config.LoadConfig(tempLog, configPath) // Pass configPath if using flags
	if err != nil {
		tempLog.Fatalf("Failed to load configuration: %v", err)
	}

	// --- Setup Real Logger ---
	log, err := logger.New(cfg.Log)
	if err != nil {
		tempLog.Fatalf("Failed to initialize logger: %v", err)
	}
	// Sync logger buffer before exiting (best effort)
	defer func() {
		if syncErr := if loggerImpl, ok := log.(*logger.zapLogger); ok { loggerImpl.sugar.Sync() }; syncErr != nil {
            // Use fmt here as logger might be failing
            fmt.Fprintf(os.Stderr, "Warning: failed to sync logger: %v\n", syncErr)
        }
	}()

	log.Infof("Configuration loaded successfully. Log Level: %s", cfg.Log.Level)

	// --- Initialize Data Sources ---
	// Instantiate only enabled sources based on cfg.DataSources
	var availableSources []dsc.DataSource
	if cfg.DataSources.WebSearch.Enable {
		// Pass only the necessary config part
		webSource, err := web.NewWebSearchSource(log, cfg.DataSources.WebSearch.APIKey)
		if err != nil {
			// Log non-fatal error and continue? Or Fatalf? Depends on requirements.
			log.Errorf("Failed to init web search source: %v. Disabling.", err)
		} else if webSource != nil {
			availableSources = append(availableSources, webSource)
		}
	}
	if cfg.DataSources.LocalFiles.Enable {
		// Pass only the necessary config part
		localSource, err := local.NewLocalFileSource(log, cfg.DataSources.LocalFiles.WatchPaths)
		if err != nil {
			log.Errorf("Failed to init local file source: %v. Disabling.", err)
		} else if localSource != nil {
			availableSources = append(availableSources, localSource)
		}
	}
	// Add other sources (DB, API) here following the same pattern...

	if len(availableSources) == 0 {
		log.Warnf("No data sources were successfully enabled or initialized!")
	}

	// --- Initialize Core Engine ---
	engine, err := qengine.NewEngine(log, cfg, availableSources)
	if err != nil {
		log.Fatalf("Failed to initialize QA engine: %v", err)
	}

	// --- Context for Graceful Shutdown ---
    // Create a context that listens for the interrupt signal from the OS.
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop() // Call stop to release resources associated with the NotifyContext

	var wg sync.WaitGroup // WaitGroup to track running servers

	// --- Setup and Run gRPC Server ---
	grpcServer := grpc.NewServer(
	// Add interceptors for logging, metrics, auth etc. here
	// grpc.UnaryInterceptor(...),
	)
	grpcAPIHandler := grpchandler.NewDataExplorerServer(engine, log)
	pb.RegisterDataExplorerServiceServer(grpcServer, grpcAPIHandler)
	reflection.Register(grpcServer) // Enable gRPC reflection

	grpcAddr := cfg.Server.GRPCPort
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen on gRPC address %s: %v", grpcAddr, err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Infof("Starting gRPC Server on %s", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			// Log error unless it's ErrServerClosed during graceful shutdown
            if err != grpc.ErrServerClosed {
			    log.Errorf("gRPC Server failed: %v", err)
                stop() // Trigger shutdown if gRPC fails critically
            }
		}
		log.Infof("gRPC Server has stopped.")
	}()

	// --- Setup and Run REST Gateway ---
	restAddr := cfg.Server.RESTPort
	wg.Add(1)
	go func() {
		defer wg.Done()
        // Pass the cancellable context `ctx` to the gateway runner
		if err := restgateway.RunRESTGateway(ctx, log, grpcAddr, restAddr); err != nil {
			log.Errorf("REST Gateway Run failed: %v", err)
            // Optional: Trigger shutdown if REST fails critically?
            // stop()
		}
        log.Infof("REST Gateway runner has finished.")
	}()


	// --- Wait for Shutdown Signal ---
	<-ctx.Done() // Block here until the context is cancelled (by signal or explicit stop call)
    stop() // Call stop() again here to ensure resources are cleaned up if shutdown was triggered externally

	log.Infof("Shutdown signal received. Initiating graceful shutdown...")

	// --- Graceful Shutdown for gRPC Server ---
	// Give it a deadline (e.g., 5 seconds)
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

    stopped := make(chan struct{})
    go func(){
        log.Infof("Attempting graceful shutdown of gRPC server...")
	    grpcServer.GracefulStop()
        close(stopped)
    }()

    // Wait for shutdown or timeout
    select {
    case <-shutdownCtx.Done():
        log.Warnf("gRPC server graceful shutdown timed out after 5s. Forcing stop.")
        grpcServer.Stop() // Force stop
    case <-stopped:
        log.Infof("gRPC server gracefully stopped.")
    }


	// Wait for all server goroutines (gRPC serve, REST run) to complete.
    // The REST gateway should stop automatically due to context cancellation.
    log.Infof("Waiting for all server goroutines to finish...")
	wg.Wait()

	log.Infof("Server shutdown complete.")
}
