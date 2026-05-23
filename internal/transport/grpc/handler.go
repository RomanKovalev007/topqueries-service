package grpctransport

import (
	"context"

	"github.com/RomanKovalev007/topqueries-service/gen/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetTopN(ctx context.Context, req *pb.GetTopRequest) (*pb.GetTopResponse, error) {
	n := int(req.GetN())
	if n <= 0 {
		return nil, status.Error(codes.InvalidArgument, "n must be greater than 0")
	}

	topQueries := s.svc.GetTopN(n)

	resp := &pb.GetTopResponse{
		Entries: make([]*pb.TopEntry, len(topQueries)),
	}
	
	for i, e := range topQueries {
		resp.Entries[i] = &pb.TopEntry{
			Query: e.Query,
			Count: e.Count,
		}
	}

	return resp, nil
}

func (s *Server) AddStopWords(ctx context.Context, req *pb.StopWordRequest) (*pb.StopWordResponse, error) {
	if len(req.GetWord()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "words list must not be empty")
	}
	// TODO: s.svc.AddStopWords(req.GetWord())
	return &pb.StopWordResponse{}, nil
}

func (s *Server) RemoveStopWords(ctx context.Context, req *pb.StopWordRequest) (*pb.StopWordResponse, error) {
	if len(req.GetWord()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "words list must not be empty")
	}
	// TODO: s.svc.RemoveStopWords(req.GetWord())
	return &pb.StopWordResponse{}, nil
}
