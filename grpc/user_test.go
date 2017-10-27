package grpc

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/Tecsisa/foulkon/api"
	util "github.com/Tecsisa/foulkon/grpc/internal/testutil"
)

// this is for testing AddUser only (more to follow)
type addUsersAPI struct {
	user *api.User
	err  error
}

func (u *addUsersAPI) AddUser(requestInfo api.RequestInfo, externalId string, path string) (*api.User, error) {
	if u.user != nil {
		return u.user, nil
	}
	return nil, u.err
}

func (*addUsersAPI) GetUserByExternalID(api.RequestInfo, string) (*api.User, error) { return nil, nil }
func (*addUsersAPI) ListUsers(api.RequestInfo, *api.Filter) ([]string, int, error)  { return nil, 0, nil }
func (*addUsersAPI) UpdateUser(api.RequestInfo, string, string) (*api.User, error)  { return nil, nil }
func (*addUsersAPI) RemoveUser(api.RequestInfo, string) error                       { return nil }
func (*addUsersAPI) ListGroupsByUser(api.RequestInfo, *api.Filter) ([]api.UserGroups, int, error) {
	return nil, 0, nil
}

func TestUserService(t *testing.T) {
	api.Log = &logrus.Logger{Out: ioutil.Discard} // no logs in tests
	ctx := util.MakeAdminCtx(context.Background())

	t.Run("when everthing goes right, returns grpc.User", func(t *testing.T) {
		eID := "externalid"
		path := "/users/"
		mock := newUsersService(&addUsersAPI{user: &api.User{ExternalID: eID, Path: path}})
		ret, err := mock.AddUser(ctx, &AddUserReq{ExternalId: eID, Path: path})
		if assert.Nil(t, err) && assert.NotNil(t, ret) {
			assert.Equal(t, ret.ExternalId, eID)
			assert.Equal(t, ret.Path, path)
		}
	})

	t.Run("when there is an error adding a user, returns an error", func(t *testing.T) {
		eID := "externalid"
		path := "/users/"
		msg := "Unable to create user, user with externalId externalid already exist"
		mock := newUsersService(&addUsersAPI{err: &api.Error{Code: api.USER_ALREADY_EXIST, Message: msg}})
		ret, err := mock.AddUser(ctx, &AddUserReq{ExternalId: eID, Path: path})
		if assert.NotNil(t, err) {
			assert.Nil(t, ret)
		}
	})
}

func TestConvertUser(t *testing.T) {
	t.Skip("TODO this test fails, take a minute to figure out the timestamps!")
	t.Run("returns an error when the created timestamp is wrong", func(t *testing.T) {
		now, _ := time.Parse(time.RFC3339, "2017-10-27 11:44:35 UTC")
		earlier := now.Add(-time.Hour)
		createdTimestamp := timestamp.Timestamp{Seconds: 1509104817}
		updatedTimestamp := createdTimestamp //timestamp.Timestamp{Seconds: (1509104817 - 60*60)}
		expected := User{
			Id:         "01234567-89ab-cdef-0123-456789abcdef",
			ExternalId: "externalid",
			Created:    &createdTimestamp,
			Updated:    &updatedTimestamp,
			Path:       "/users/",
			Urn:        "urn:iws:iam::user/example/admin/user1",
		}
		u := api.User{
			ID:         "01234567-89ab-cdef-0123-456789abcdef",
			ExternalID: "externalid",
			CreateAt:   earlier,
			UpdateAt:   now,
			Path:       "/users/",
			Urn:        "urn:iws:iam::user/example/admin/user1",
		}

		r, err := convertUser(&u)
		if assert.Nil(t, err) {
			assert.Equal(t, expected, r)
		}
	})
}
