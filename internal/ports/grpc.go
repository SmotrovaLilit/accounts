package ports

import (
	"accounts/internal/app"
	pb "accounts/internal/gen-proto/accountservice/v1beta"
)

// GRPCServer is a gRPC server
type GRPCServer struct {
	app app.Application
	pb.UnimplementedAccountServiceServer
}

// NewGRPCServer creates a new gRPC server
func NewGRPCServer(application app.Application) GRPCServer {
	return GRPCServer{
		app: application,
	}
}
