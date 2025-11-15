package grpc

import (
	"context"
	pb "github.com/mikhailbolshakov/ocpi/proto"
)

func (s *Server) PullParties(ctx context.Context, rq *pb.EmptyRequest) (*pb.EmptyResponse, error) {
	err := s.credUc.OnRemotePartyPull(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}

func (s *Server) PullHubClientInfo(ctx context.Context, rq *pb.RemotePullRequest) (*pb.EmptyResponse, error) {
	from, to := s.toPeriodUc(rq)
	err := s.hubUc.OnRemoteClientInfosPull(ctx, from, to)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}

func (s *Server) PullHubClientInfoWhenPushNotSupported(ctx context.Context, rq *pb.RemotePullRequest) (*pb.EmptyResponse, error) {
	from, to := s.toPeriodUc(rq)
	err := s.hubUc.OnRemoteClientInfosPullWhenPushNotSupported(ctx, from, to)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}

func (s *Server) PullLocations(ctx context.Context, rq *pb.RemotePullRequest) (*pb.EmptyResponse, error) {
	from, to := s.toPeriodUc(rq)
	err := s.locUc.OnRemoteLocationsPull(ctx, from, to)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}

func (s *Server) PullLocationsWhenPushNotSupported(ctx context.Context, rq *pb.RemotePullRequest) (*pb.EmptyResponse, error) {
	from, to := s.toPeriodUc(rq)
	err := s.locUc.OnRemoteLocationsPullWhenPushNotSupported(ctx, from, to)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}

func (s *Server) PullTariffs(ctx context.Context, rq *pb.RemotePullRequest) (*pb.EmptyResponse, error) {
	from, to := s.toPeriodUc(rq)
	err := s.tariffUc.OnRemoteTariffsPull(ctx, from, to)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}

func (s *Server) PullTariffsWhenPushNotSupported(ctx context.Context, rq *pb.RemotePullRequest) (*pb.EmptyResponse, error) {
	from, to := s.toPeriodUc(rq)
	err := s.tariffUc.OnRemoteTariffsPullWhenPushNotSupported(ctx, from, to)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}

func (s *Server) PullTokens(ctx context.Context, rq *pb.RemotePullRequest) (*pb.EmptyResponse, error) {
	from, to := s.toPeriodUc(rq)
	err := s.tokenUc.OnRemoteTokensPull(ctx, from, to)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}

func (s *Server) PullTokensWhenPushNotSupported(ctx context.Context, rq *pb.RemotePullRequest) (*pb.EmptyResponse, error) {
	from, to := s.toPeriodUc(rq)
	err := s.tokenUc.OnRemoteTokensPullWhenPushNotSupported(ctx, from, to)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}

func (s *Server) PullSessions(ctx context.Context, rq *pb.RemotePullRequest) (*pb.EmptyResponse, error) {
	from, to := s.toPeriodUc(rq)
	err := s.sessUc.OnRemoteSessionsPull(ctx, from, to)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}

func (s *Server) PullSessionsWhenPushNotSupported(ctx context.Context, rq *pb.RemotePullRequest) (*pb.EmptyResponse, error) {
	from, to := s.toPeriodUc(rq)
	err := s.sessUc.OnRemoteSessionsPullWhenPushNotSupported(ctx, from, to)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}

func (s *Server) PullCdrs(ctx context.Context, rq *pb.RemotePullRequest) (*pb.EmptyResponse, error) {
	from, to := s.toPeriodUc(rq)
	err := s.cdrUc.OnRemoteCdrsPull(ctx, from, to)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}

func (s *Server) PullCdrsWhenPushNotSupported(ctx context.Context, rq *pb.RemotePullRequest) (*pb.EmptyResponse, error) {
	from, to := s.toPeriodUc(rq)
	err := s.cdrUc.OnRemoteCdrsPullWhenPushNotSupported(ctx, from, to)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}
