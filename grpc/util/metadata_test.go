package util_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/grpc/util"
)

func TestRequestInfoPassedThroughContext(t *testing.T) {
	reqInfo := &api.RequestInfo{
		Admin:      true,
		Identifier: "identifier",
		RequestID:  "asdfghjkl",
	}
	ctx := util.ContextWithRequestInfo(context.Background(), reqInfo)
	reqInfo2, ok := util.RequestInfoFromContext(ctx)
	if assert.True(t, ok) {
		assert.Equal(t, reqInfo, reqInfo2)
	}
}
