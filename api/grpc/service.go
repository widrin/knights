package grpc

// GRPCServer provides gRPC API endpoints
type GRPCServer struct {
	address string
}

// NewGRPCServer creates a new gRPC server
func NewGRPCServer(address string) *GRPCServer {
	return &GRPCServer{
		address: address,
	}
}

// Start starts the gRPC server
func (s *GRPCServer) Start() error {
	// TODO: Implement gRPC server
	return nil
}
