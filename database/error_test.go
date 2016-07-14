package database

import (
	"github.com/kylelemons/godebug/pretty"
	"testing"
)

func TestError(t *testing.T) {
	testcases := map[string]struct {
		error           Error
		expectedMessage string
	}{
		"EmptyValues": {
			error:           Error{},
			expectedMessage: "Code: , Message: ",
		},
		"EmptyCode": {
			error: Error{
				Message: "This is a message",
			},
			expectedMessage: "Code: , Message: This is a message",
		},
		"EmptyMessage": {
			error: Error{
				Code: "CODE",
			},
			expectedMessage: "Code: CODE, Message: ",
		},
		"NormalCase": {
			error: Error{
				Code:    INTERNAL_ERROR,
				Message: "Internal error",
			},
			expectedMessage: "Code: " + INTERNAL_ERROR + ", Message: Internal error",
		},
	}

	for x, testcase := range testcases {
		if diff := pretty.Compare(testcase.expectedMessage, testcase.error.Error()); diff != "" {
			t.Errorf("Test %v failed. Received different messages (wanted / received) %v",
				x, diff)
			continue
		}
	}
}
