package grpctransport

import (
	"github.com/RomanKovalev007/topqueries-service/gen/pb"
	"github.com/RomanKovalev007/topqueries-service/internal/domain"
)

type service interface{
	Add(query string)
	GetTopN(n int) []domain.TopEntry
}

type Server struct{
	pb.UnimplementedTopQueriesServiceServer
	svc service
}

func NewServer(svc service) *Server {
	return &Server{ svc: svc }
}