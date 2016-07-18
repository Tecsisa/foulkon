package api

import (
	"bytes"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/kylelemons/godebug/pretty"
	"strings"
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
		if testcase.err.Error() != testcase.expectedMessage {
			t.Errorf("Test %v failed. Received different messages (wanted:%v / received:%v)",
				x, testcase.expectedMessage, testcase.err.Error())
			continue
		}
	}
}

func TestLogErrorMessage(t *testing.T) {
	// Logger
	testOut := bytes.NewBuffer([]byte{})
	logger := logrus.Logger{
		Out:       testOut,
		Formatter: &logrus.JSONFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.ErrorLevel,
	}
	testcases := map[string]struct {
		err             *Error
		expectedMessage string
	}{
		"OkCase": {
			err: &Error{
				Code:    "Code Error",
				Message: "Error Message",
			},
			expectedMessage: "{\"Code\":\"Code Error\",\"RequestID\":\"RequestID\",\"level\":\"error\",\"msg\":\"Error Message\",\"time\"",
		},
	}

	for x, testcase := range testcases {
		LogErrorMessage(&logger, "RequestID", testcase.err)
		logMessage := testOut.String()
		diff := pretty.Compare(logMessage, testcase.expectedMessage)
		if !strings.Contains(logMessage, testcase.expectedMessage) {
			t.Errorf("Test %v failed. Received different messages (wanted / received) %v",
				x, diff)
			continue
		}
	}
}
