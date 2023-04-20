package grpc

import (
	"infoblox-golang/cmd/addressbook/transport/grpc/pb"
	"infoblox-golang/internal/addressbook"
	"infoblox-golang/internal/platform/storage"
)

func UserToProto(u *addressbook.User) *pb.User {
	return &pb.User{
		Id:       string(u.ID),
		Username: u.Username,
		Address:  u.Address,
		Phone:    u.Phone,
	}
}

func ProtoToUser(p *pb.User) *addressbook.User {
	return &addressbook.User{
		ID:       storage.ID(p.Id),
		Username: p.Username,
		Address:  p.Address,
		Phone:    p.Phone,
	}
}
