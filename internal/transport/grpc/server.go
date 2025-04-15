package grpc

import (
	"context"
	"fmt" // Added for error formatting

	pb "github.com/yourusername/data-explorer/api/proto/dataexplorer/v1" // Adjust path
	"github.com/yourusername/data-explorer/internal/domain"       // Adjust path
	"github.com/yourusername/data-explorer/internal/platform/logger"    // Adjust path
	"github.com/yourusername/data-explorer/internal/qengine"      // Adjust path
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	// Consider using "google.golang.org/genproto/googleapis/rpc/errdetails" for richer errors
)

// DataExplorerServer implements the gRPC service generated from the proto definition.
type DataExplorerServer struct {
	// Embed the unimplemented server type for forward compatibility.
	// Ensures your server implementation stays compatible if new methods are added to the service definition.
	pb.UnimplementedDataExplorerServiceServer
	engine qengine.Engine
	log    logger.Logger
}

// NewDataExplorerServer creates a new gRPC server handler instance.
func NewDataExplorerServer(engine qengine.Engine, log logger.Logger) *DataExplorerServer {
	return &DataExplorerServer{
		engine: engine,
		log:    log.With("transport", "grpc"), // Add structured context to logs
	}
}

// AskQuestion handles the AskQuestion gRPC request.
func (s *DataExplorerServer) AskQuestion(ctx context.Context, req *pb.AskQuestionRequest) (*pb.AskQuestionResponse, error) {
	// Log the incoming request (avoid logging sensitive data in production if query can be sensitive)
	s.log.Infof("Received AskQuestion request: Query length=%d, NumSources=%d, NumOptions=%d",
		len(req.Query), len(req.Sources), len(req.Options))
    s.log.Debugf("Full Query (debug): %s", req.Query) // Log full query only at debug level

	// --- Input Validation ---
	if req.Query == "" {
		s.log.Warnf("Received empty query in AskQuestion request")
		// Return a gRPC status error with InvalidArgument code.
		return nil, status.Error(codes.InvalidArgument, "query cannot be empty")
	}
	// TODO: Add more validation for req.Sources and req.Options if needed.

	// --- Map Protobuf Request to Internal Domain Model ---
	domainQuestion := domain.Question{
		Query:   req.Query,
		Options: req.Options,
		Sources: make([]domain.DataSourceSpec, len(req.Sources)),
	}
	for i, src := range req.Sources {
		domainQuestion.Sources[i] = domain.DataSourceSpec{
			Type:       src.Type,
			Identifier: src.Identifier,
		}
	}

	// --- Call the Core Logic (QA Engine) ---
	// The context passed from gRPC handles deadlines and cancellations.
	domainAnswer, err := s.engine.ProcessQuestion(ctx, domainQuestion)

	// --- Handle Engine Errors ---
	if err != nil {
		s.log.Errorf("Error processing question in engine: %v", err)
		// TODO: Map specific internal errors to appropriate gRPC status codes.
		// For now, return a generic Internal error for engine failures.
		// Consider checking error type (e.g., using errors.Is or type assertion)
		// to return codes like NotFound, PermissionDenied, etc.
		return nil, status.Errorf(codes.Internal, "failed to process question: %v", err)
	}
    if domainAnswer == nil {
        // This case should ideally not happen if the engine guarantees a non-nil answer or an error.
        s.log.Errorf("Engine returned nil answer without error, indicating an internal issue.")
        return nil, status.Error(codes.Internal, "internal error: engine returned unexpected nil answer")
    }

	// --- Map Internal Domain Model back to Protobuf Response ---
	resp := &pb.AskQuestionResponse{
		Answer:      domainAnswer.Text,
		SourcesUsed: make([]*pb.SourceDocument, len(domainAnswer.SourcesUsed)),
	}
	// Include non-fatal processing errors/warnings from the engine in the response payload.
	if domainAnswer.Error != nil {
		s.log.Warnf("Including non-fatal error in response: %v", domainAnswer.Error)
		resp.Error = domainAnswer.Error.Error() // Simple string representation for now
	}
	// Map source documents
	for i, srcDoc := range domainAnswer.SourcesUsed {
		resp.SourcesUsed[i] = &pb.SourceDocument{
			SourceType:     srcDoc.SourceType,
			Identifier:     srcDoc.Identifier,
			ContentSnippet: srcDoc.ContentSnippet,
			RelevanceScore: srcDoc.RelevanceScore,
		}
	}

	s.log.Infof("Successfully processed AskQuestion request for query length %d", len(req.Query))
	// Return the response and a nil error to indicate success to gRPC.
	return resp, nil
}
