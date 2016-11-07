package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	testcases := map[string]struct {
		err             Error
		expectedMessage string
	}{
		"EmptyValues": {
			err:             Error{},
			expectedMessage: "Code: , Message: ",
		},
		"EmptyCode": {
			err: Error{
				Message: "This is a message",
			},
			expectedMessage: "Code: , Message: This is a message",
		},
		"EmptyMessage": {
			err: Error{
				Code: "CODE",
			},
			expectedMessage: "Code: CODE, Message: ",
		},
		"NormalCase": {
			err: Error{
				Code:    UNKNOWN_API_ERROR,
				Message: "Unknown error",
			},
			expectedMessage: "Code: " + UNKNOWN_API_ERROR + ", Message: Unknown error",
		},
	}

	for x, testcase := range testcases {
		assert.Equal(t, testcase.expectedMessage, testcase.err.Error(), "Error in test case %v", x)
	}
}
