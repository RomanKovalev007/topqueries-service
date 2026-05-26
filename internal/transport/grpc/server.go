package grpctransport

import (
	"github.com/RomanKovalev007/topqueries-service/gen/pb"
	"github.com/RomanKovalev007/topqueries-service/internal/domain"
)

type service interface{
	GetTopN(n int) []domain.TopEntry
	AddStopWords(words []string)
	RemoveStopWords(words []string)
}

type Server struct{
	pb.UnimplementedTopQueriesServiceServer
	svc service
}

func NewServer(svc service) *Server {
	return &Server{ svc: svc }
}