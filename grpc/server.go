package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/foulkon"
)

// Server struct
type Server struct {
	certFile string
	keyFile  string
	addr     string
	userAPI  api.UserAPI

	s *grpc.Server
}

// NewServer returns a new Server
func NewServer(worker *foulkon.Worker) *Server {
	return &Server{
		certFile: worker.CertFile,
		keyFile:  worker.KeyFile,
		addr:     worker.GRPCHost + ":" + worker.GRPCPort,
		userAPI:  worker.UserApi,
	}
}

// Run starts a GRPCServer
func (gs *Server) Run() error {
	// optionally set up TLS
	if gs.certFile != "" && gs.keyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(gs.certFile, gs.keyFile)
		if err != nil {
			return err
		}
		gs.s = grpc.NewServer(grpc.Creds(creds))
	} else {
		gs.s = grpc.NewServer()
	}

	// Register reflection service (to lookup exposed services)
	reflection.Register(gs.s)

	// register our services
	RegisterUserServiceServer(gs.s, newUsersService(gs.userAPI))

	lis, err := net.Listen("tcp", gs.addr)
	if err != nil {
		return fmt.Errorf("listen on %s failed: %s", gs.addr, err)
	}
	return gs.s.Serve(lis)
}
