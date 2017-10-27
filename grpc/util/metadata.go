package util

import (
	"context"

	"github.com/Tecsisa/foulkon/api"
)

type reqInfoKey struct{}

// ContextWithRequestInfo injects the passed RequestInfo into the passed
// context. It can be retrieved with RequestInfoFromContext.
func ContextWithRequestInfo(ctx context.Context, info *api.RequestInfo) context.Context {
	return context.WithValue(ctx, reqInfoKey{}, info)
}

// RequestInfoFromContext tries to extract RequestInfo from the passed context.
// On success, returns the info and true; otherwise false.
func RequestInfoFromContext(ctx context.Context) (*api.RequestInfo, bool) {
	ri, ok := ctx.Value(reqInfoKey{}).(*api.RequestInfo)
	return ri, ok
}
