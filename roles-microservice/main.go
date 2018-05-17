package main

import (
	"errors"
	pb "github.com/rymccue/grpc-communication-demo/roles-microservice/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	userRoles map[int32][]*pb.Role
	roles     []*pb.Role
}

func (s *Server) GetRoles(_ context.Context, _ *pb.EmptyRequest) (*pb.RolesReply, error) {
	return &pb.RolesReply{
		Roles: s.roles,
	}, nil
}

func (s *Server) GetUserRole(_ context.Context, req *pb.GetUserRoleRequest) (*pb.UserRoleReply, error) {
	roles, ok := s.userRoles[req.UserId]
	if !ok {
		return nil, errors.New("user not found")
	}
	return &pb.UserRoleReply{
		UserId: req.UserId,
		Roles:  roles,
	}, nil
}

func main() {

	var (
		normal = &pb.Role{
			Id:   1,
			Name: "normal",
		}
		editor = &pb.Role{
			Id:   2,
			Name: "editor",
		}
		admin = &pb.Role{
			Id:   3,
			Name: "admin",
		}
		superUser = &pb.Role{
			Id:   4,
			Name: "super user",
		}
	)

	lis, err := net.Listen("tcp", "localhost:6000")
	if err != nil {
		log.Fatalf("failed to initializa TCP listen: %v", err)
	}
	defer lis.Close()

	server := grpc.NewServer()
	roleServer := &Server{
		userRoles: map[int32][]*pb.Role{
			1: {normal},
			2: {normal, editor},
			3: {normal},
			4: {normal, editor, admin},
			5: {normal, editor, admin, superUser},
		},
		roles: []*pb.Role{normal, editor, admin, superUser},
	}
	pb.RegisterRolesServer(server, roleServer)

	server.Serve(lis)
}
