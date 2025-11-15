package grpc

import (
	kitGrpc "github.com/mikhailbolshakov/kit/grpc"
	"github.com/mikhailbolshakov/ocpi"
	pb "github.com/mikhailbolshakov/ocpi/proto"
	"github.com/mikhailbolshakov/ocpi/usecase"
)

// Server implements gRPC server
type Server struct {
	*kitGrpc.Server

	credUc   usecase.CredentialsUc
	locUc    usecase.LocationUc
	hubUc    usecase.HubUc
	tariffUc usecase.TariffUc
	tokenUc  usecase.TokenUc
	sessUc   usecase.SessionUc
	cdrUc    usecase.CdrUc

	pb.UnimplementedRemotePullServiceServer
}

// New creates a new gRPC server
func New(
	credUc usecase.CredentialsUc,
	locUc usecase.LocationUc,
	hubUc usecase.HubUc,
	tariffUc usecase.TariffUc,
	tokenUc usecase.TokenUc,
	sessUc usecase.SessionUc,
	cdrUc usecase.CdrUc) *Server {
	return &Server{
		credUc:   credUc,
		locUc:    locUc,
		hubUc:    hubUc,
		tariffUc: tariffUc,
		tokenUc:  tokenUc,
		sessUc:   sessUc,
		cdrUc:    cdrUc,
	}
}

// Init initializes gRPC server
func (s *Server) Init(cfg *kitGrpc.ServerConfig) error {
	// grpc server
	gs, err := kitGrpc.NewServer(ocpi.Meta.ServiceCode(), ocpi.LF(), cfg)
	if err != nil {
		return err
	}
	s.Server = gs
	pb.RegisterRemotePullServiceServer(s.Srv, s)
	return nil
}
