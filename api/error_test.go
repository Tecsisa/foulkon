package api

import "testing"

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
				Code:    UNKNOWN_API_ERROR,
				Message: "Unknown error",
			},
			expectedMessage: "Code: " + UNKNOWN_API_ERROR + ", Message: Unknown error",
		},
	}

	for x, testcase := range testcases {
		if testcase.error.Error() != testcase.expectedMessage {
			t.Errorf("Test %v failed. Received different messages (wanted:%v / received:%v)",
				x, testcase.expectedMessage, testcase.error.Error())
			continue
		}
	}
}
