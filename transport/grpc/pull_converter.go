package grpc

import (
	"github.com/mikhailbolshakov/kit/grpc"
	pb "github.com/mikhailbolshakov/ocpi/proto"
	"time"
)

func (s *Server) toPeriodUc(rq *pb.RemotePullRequest) (*time.Time, *time.Time) {
	return grpc.PbTSToTime(rq.From), grpc.PbTSToTime(rq.To)
}
