package rest

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/yourusername/data-explorer/api/proto/dataexplorer/v1" // Adjust
	"github.com/yourusername/data-explorer/internal/platform/logger"    // Adjust
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // Use secure creds in production!
)

// RunRESTGateway starts the REST gateway proxy for the gRPC service.
func RunRESTGateway(ctx context.Context, log logger.Logger, grpcServerEndpoint, restListenAddr string) error {
	log = log.With("transport", "rest_gateway")
	log.Infof("Starting REST Gateway proxy...")

	// Create a new ServeMux for the gateway
	// Use runtime.WithErrorHandler for custom error handling if needed
	// Use runtime.WithMarshalerOption for customizing JSON marshaling
	mux := runtime.NewServeMux()

	// Define gRPC client connection options
	// IMPORTANT: Use secure credentials (TLS) in production environments!
	// opts := []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())} // For local dev/testing

	// Register the gRPC service handler endpoint
	log.Infof("Registering gRPC handler endpoint: %s -> REST Gateway", grpcServerEndpoint)
	err := pb.RegisterDataExplorerServiceHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts)
	if err != nil {
		log.Errorf("Failed to register gRPC Gateway handler: %v", err)
		return err
	}

	// Create and start the HTTP server for the REST Gateway
	log.Infof("REST Gateway listening on %s", restListenAddr)
	httpServer := &http.Server{
		Addr:    restListenAddr,
		Handler: corsMiddleware(mux), // Add CORS middleware if needed
	}

	// Run the server in a goroutine so it doesn't block.
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("REST Gateway failed to listen: %v", err) // Use Fatalf to exit if listen fails critically
		}
	}()

    // Listen for the context being cancelled (e.g., shutdown signal)
    <-ctx.Done()

    // Context cancelled, initiate graceful shutdown
    log.Infof("Shutting down REST Gateway server...")

    // Create a deadline context for the shutdown
    // shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Use time package
    // defer cancel() // Ensure cancel is called eventually

    // Attempt graceful shutdown
    // if err := httpServer.Shutdown(shutdownCtx); err != nil {
    //     log.Errorf("REST Gateway graceful shutdown failed: %v", err)
    //     return err // Or handle more gracefully
    // }

    // Temporary simpler shutdown (replace with above for production)
    if err := httpServer.Close(); err != nil {
        log.Errorf("REST Gateway close failed: %v", err)
        return err
    }


	log.Infof("REST Gateway server stopped gracefully.")
	return nil
}

// corsMiddleware is a basic example of adding CORS headers. Use a robust library for production.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin for simplicity (restrict in production!)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web") // Add necessary headers
        w.Header().Set("Access-Control-Expose-Headers", "Grpc-Metadata-*, X-Grpc-Metadata-*") // Expose gRPC metadata headers if needed

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
