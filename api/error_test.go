package api

import "testing"

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
		if testcase.err.Error() != testcase.expectedMessage {
			t.Errorf("Test %v failed. Received different messages (wanted:%v / received:%v)",
				x, testcase.expectedMessage, testcase.err.Error())
			continue
		}
	}
}
