package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/grpc/util"
)

type users struct {
	api api.UserAPI
}

func newUsersService(api api.UserAPI) *users {
	return &users{api: api}
}

func (u *users) AddUser(ctx context.Context, req *AddUserReq) (*User, error) {
	// TODO reqInfoFromContext, add it into ctx via grpc middleware
	reqInfo, ok := util.RequestInfoFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, api.AUTHENTICATION_API_ERROR)
	}
	// For manual testing purposes, comment out the call above and comment in the line below
	// reqInfo := &api.RequestInfo{Admin: true}
	user, err := u.api.AddUser(*reqInfo, req.ExternalId, req.Path)
	if err, ok := err.(*api.Error); ok {
		api.LogOperationError(reqInfo.RequestID, reqInfo.Identifier, err)
		return nil, statusFromError(err)
	}

	ret, err := convertUser(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process user response: %s", err)
	}
	return ret, nil
}

// convert api.User to grpc.User
func convertUser(user *api.User) (*User, error) {
	created, err := ptypes.TimestampProto(user.CreateAt)
	if err != nil {
		return nil, err
	}
	updated, err := ptypes.TimestampProto(user.UpdateAt)
	if err != nil {
		return nil, err
	}
	return &User{
		ExternalId: user.ExternalID,
		Path:       user.Path,
		Id:         user.ID,
		Urn:        user.GetUrn(),
		Created:    created,
		Updated:    updated,
	}, nil
}
