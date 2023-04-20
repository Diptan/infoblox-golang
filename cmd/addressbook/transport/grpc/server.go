package grpc

import (
	"context"
	"fmt"
	"infoblox-golang/cmd/addressbook/transport/grpc/pb"
	"infoblox-golang/internal/addressbook"
	"infoblox-golang/internal/platform/utils"
	"net"
	"net/http"
	"reflect"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedAddressBookServer
	service *addressbook.Service
	srv     *grpc.Server
}

func NewServer(service *addressbook.Service) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) CreateUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	us := ProtoToUser(user)

	u, err := s.service.AddUser(*us)
	if err != nil {
		status.Errorf(codes.Internal, "Fail to create user: %s", err)
	}
	uPb := UserToProto(&u)

	return uPb, nil
}

func (s *Server) DeleteUser(ctx context.Context, id *pb.UserId) (*pb.DeleteUserResponse, error) {

	err := s.service.DeleteUser(id.Id)
	if err != nil {
		status.Errorf(codes.Internal, "Fail to delete user: %s", err)
	}

	return &pb.DeleteUserResponse{Success: true}, nil
}

func (s *Server) UpdateUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	_, err := s.service.UpdateUser(user.Id, *ProtoToUser(user))
	if err != nil {
		status.Errorf(codes.Internal, "Fail to update user: %s", err)
	}

	return user, nil
}

func (s *Server) GetAll(context.Context, *emptypb.Empty) (*pb.Users, error) {
	users, err := s.service.GetAllUsers()
	if err != nil {
		status.Errorf(codes.Internal, "Fail to get users: %s", err)
	}

	var pbUsers pb.Users
	for _, u := range users {
		pbUsers.Users = append(pbUsers.Users, UserToProto(&u))
	}

	return &pbUsers, nil
}

// TODO: Refactor this method -> move isMatch logic to service layer
func (s *Server) SearchUser(ctx context.Context, req *pb.SearchUserRequest) (*pb.Users, error) {
	// Define a slice to hold the matching users
	users, err := s.service.GetAllUsers()
	if err != nil {
		status.Errorf(codes.Internal, "Fail to get users: %s", err)
	}

	var matchingUsers []*pb.User

	// Loop through each existing user
	for _, user := range users {

		//Try to get pb.User struct field by the string provided in request
		field := reflect.ValueOf(user).FieldByNameFunc(func(n string) bool { return strings.ToLower(n) == req.Field })
		if !field.IsValid() {
			return nil, status.Errorf(codes.InvalidArgument, "Invalid field value: %v", req.Field)
		}

		if utils.IsMatch(field.String(), req.Value) {
			matchingUsers = append(matchingUsers, UserToProto(&user))
		}
	}

	if len(matchingUsers) == 0 {
		return nil, status.Errorf(codes.NotFound, "Records not found. Field: %s value: %s", req.Field, req.Value)
	}

	return &pb.Users{
		Users: matchingUsers,
	}, nil
}

// Start simple grpc server (This method was implemented for education purposes)
func (s *Server) ListenAndServe(port int) error {
	addr := fmt.Sprintf(":%d", port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	err = s.srv.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

// Start HTTP server (and proxy calls to gRPC server endpoint)
func (s *Server) RunGatewayServer(ctx context.Context, port int) error {

	grpcMux := runtime.NewServeMux()

	err := pb.RegisterAddressBookHandlerServer(ctx, grpcMux, s)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	addr := fmt.Sprintf(":%d", port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	err = http.Serve(lis, mux)
	if err != nil {
		return err
	}

	return nil
}
