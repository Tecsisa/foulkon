package testutil

import (
	"context"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/grpc/util"
)

// MakeAdminCtx returns a context whose derived RequestInfo indicates an admin
// request
func MakeAdminCtx(ctx context.Context) context.Context {
	reqInfo := api.RequestInfo{
		Admin:      true,
		Identifier: "",
		RequestID:  "",
	}
	return util.ContextWithRequestInfo(ctx, &reqInfo)
}
