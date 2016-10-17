package api

import (
	"testing"

	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestLogOperation(t *testing.T) {
	testcases := map[string]struct {
		requestID string
		userID    string
		message   string
	}{
		"OkTestCase": {
			requestID: "123",
			userID:    "user123",
			message:   "message",
		},
	}

	testLogger, hook := test.NewNullLogger()
	Log = testLogger
	for n, testcase := range testcases {
		LogOperation(testcase.requestID, testcase.userID, testcase.message)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level, "Error in test case %v", n)
		assert.Equal(t, testcase.message, hook.LastEntry().Message, "Error in test case %v", n)
		assert.Equal(t, testcase.requestID, hook.LastEntry().Data["requestID"], "Error in test case %v", n)
		assert.Equal(t, testcase.userID, hook.LastEntry().Data["user"], "Error in test case %v", n)
	}
}

func TestLogOperationError(t *testing.T) {
	testcases := map[string]struct {
		requestID string
		userID    string
		err       *Error
	}{
		"OkTestCase": {
			requestID: "123",
			userID:    "user123",
			err: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Message error",
			},
		},
	}

	testLogger, hook := test.NewNullLogger()
	Log = testLogger
	for n, testcase := range testcases {
		LogOperationError(testcase.requestID, testcase.userID, testcase.err)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level, "Error in test case %v", n)
		assert.Equal(t, testcase.err.Message, hook.LastEntry().Message, "Error in test case %v", n)
		assert.Equal(t, testcase.requestID, hook.LastEntry().Data["requestID"], "Error in test case %v", n)
		assert.Equal(t, testcase.userID, hook.LastEntry().Data["user"], "Error in test case %v", n)
		assert.Equal(t, testcase.err.Code, hook.LastEntry().Data["errorCode"], "Error in test case %v", n)
	}
}

func TestTransactionRequestLog(t *testing.T) {
	testcases := map[string]struct {
		requestID   string
		userID      string
		httpMethod  string
		httpURI     string
		httpAddress string
	}{
		"OkTestCase": {
			requestID:   "123",
			userID:      "user123",
			httpMethod:  http.MethodGet,
			httpURI:     "/get",
			httpAddress: "localhost",
		},
	}

	testLogger, hook := test.NewNullLogger()
	Log = testLogger
	for n, testcase := range testcases {
		req, err := http.NewRequest(testcase.httpMethod, testcase.httpAddress+testcase.httpURI, nil)
		if err != nil {
			t.Errorf("Test case %v. Received unexpected error creating request: %v", testcase, err)
		}
		TransactionRequestLog(testcase.requestID, testcase.userID, req)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level, "Error in test case %v", n)
		assert.Equal(t, "", hook.LastEntry().Message, "Error in test case %v", n)
		assert.Equal(t, testcase.requestID, hook.LastEntry().Data["requestID"], "Error in test case %v", n)
		assert.Equal(t, testcase.httpMethod, hook.LastEntry().Data["httpMethod"], "Error in test case %v", n)
		assert.Equal(t, testcase.httpAddress+testcase.httpURI, hook.LastEntry().Data["httpURI"], "Error in test case %v", n)
		assert.Equal(t, "", hook.LastEntry().Data["httpRemoteAddress"], "Error in test case %v", n)
	}
}

func TestTransactionResponseErrorLog(t *testing.T) {
	testcases := map[string]struct {
		requestID      string
		err            *Error
		httpMethod     string
		httpURI        string
		httpAddress    string
		httpStatusCode int
	}{
		"OkTestCase": {
			requestID: "123",
			err: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Message error",
			},
			httpMethod:     http.MethodGet,
			httpURI:        "/get",
			httpAddress:    "localhost",
			httpStatusCode: http.StatusOK,
		},
	}

	testLogger, hook := test.NewNullLogger()
	Log = testLogger
	for n, testcase := range testcases {
		req, err := http.NewRequest(testcase.httpMethod, testcase.httpAddress+testcase.httpURI, nil)
		if err != nil {
			t.Errorf("Test case %v. Received unexpected error creating request: %v", testcase, err)
		}
		TransactionResponseErrorLog(testcase.requestID, "", req, testcase.httpStatusCode, testcase.err)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level, "Error in test case %v", n)
		assert.Equal(t, testcase.err.Message, hook.LastEntry().Message, "Error in test case %v", n)
		assert.Equal(t, testcase.requestID, hook.LastEntry().Data["requestID"], "Error in test case %v", n)
		assert.Equal(t, testcase.err.Code, hook.LastEntry().Data["errorCode"], "Error in test case %v", n)
		assert.Equal(t, testcase.httpMethod, hook.LastEntry().Data["httpMethod"], "Error in test case %v", n)
		assert.Equal(t, testcase.httpAddress+testcase.httpURI, hook.LastEntry().Data["httpURI"], "Error in test case %v", n)
		assert.Equal(t, "", hook.LastEntry().Data["httpRemoteAddress"], "Error in test case %v", n)
		assert.Equal(t, testcase.httpStatusCode, hook.LastEntry().Data["httpStatusCode"], "Error in test case %v", n)
		hook.Reset()
	}
}

func TestTransactionProxyLog(t *testing.T) {
	testcases := map[string]struct {
		requestID       string
		workerRequestId string
		message         string
		httpMethod      string
		httpURI         string
		httpAddress     string
	}{
		"OkTestCase": {
			requestID:       "123",
			workerRequestId: "456",
			message:         "Message",
			httpMethod:      http.MethodGet,
			httpURI:         "/get",
			httpAddress:     "localhost",
		},
		"OkTestCaseEmptyWorkerRequestId": {
			requestID:       "123",
			workerRequestId: "",
			message:         "Message",
			httpMethod:      http.MethodPost,
			httpURI:         "/post",
			httpAddress:     "localhost",
		},
	}

	testLogger, hook := test.NewNullLogger()
	Log = testLogger
	for n, testcase := range testcases {
		req, err := http.NewRequest(testcase.httpMethod, testcase.httpAddress+testcase.httpURI, nil)
		if err != nil {
			t.Errorf("Test case %v. Received unexpected error creating request: %v", testcase, err)
		}
		TransactionProxyLog(testcase.requestID, testcase.workerRequestId, req, testcase.message)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level, "Error in test case %v", n)
		assert.Equal(t, testcase.message, hook.LastEntry().Message, "Error in test case %v", n)
		assert.Equal(t, testcase.requestID, hook.LastEntry().Data["requestID"], "Error in test case %v", n)
		if testcase.workerRequestId == "" {
			assert.Empty(t, hook.LastEntry().Data["workerRequestID"], "Error in test case %v", n)
		} else {
			assert.Equal(t, testcase.workerRequestId, hook.LastEntry().Data["workerRequestID"], "Error in test case %v", n)
		}
		assert.Equal(t, testcase.httpMethod, hook.LastEntry().Data["httpMethod"], "Error in test case %v", n)
		assert.Equal(t, testcase.httpAddress+testcase.httpURI, hook.LastEntry().Data["httpURI"], "Error in test case %v", n)
		assert.Equal(t, "", hook.LastEntry().Data["httpRemoteAddress"], "Error in test case %v", n)
		hook.Reset()
	}
}

func TestTransactionProxyErrorLogWithStatus(t *testing.T) {
	testcases := map[string]struct {
		requestID       string
		workerRequestId string
		err             *Error
		httpMethod      string
		httpURI         string
		httpAddress     string
		httpStatusCode  int
	}{
		"OkTestCase": {
			requestID:       "123",
			workerRequestId: "456",
			err: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Message error",
			},
			httpMethod:     http.MethodGet,
			httpURI:        "/get",
			httpAddress:    "localhost",
			httpStatusCode: http.StatusOK,
		},
		"OkTestCaseEmptyWorkerRequestId": {
			requestID:       "123",
			workerRequestId: "",
			err: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Message error",
			},
			httpMethod:     http.MethodPost,
			httpURI:        "/post",
			httpAddress:    "localhost",
			httpStatusCode: http.StatusInternalServerError,
		},
	}

	testLogger, hook := test.NewNullLogger()
	Log = testLogger
	for n, testcase := range testcases {
		req, err := http.NewRequest(testcase.httpMethod, testcase.httpAddress+testcase.httpURI, nil)
		if err != nil {
			t.Errorf("Test case %v. Received unexpected error creating request: %v", testcase, err)
		}
		TransactionProxyErrorLogWithStatus(testcase.requestID, testcase.workerRequestId, req, testcase.httpStatusCode, testcase.err)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level, "Error in test case %v", n)
		assert.Equal(t, testcase.err.Message, hook.LastEntry().Message, "Error in test case %v", n)
		assert.Equal(t, testcase.requestID, hook.LastEntry().Data["requestID"], "Error in test case %v", n)
		if testcase.workerRequestId == "" {
			assert.Empty(t, hook.LastEntry().Data["workerRequestID"], "Error in test case %v", n)
		} else {
			assert.Equal(t, testcase.workerRequestId, hook.LastEntry().Data["workerRequestID"], "Error in test case %v", n)
		}
		assert.Equal(t, testcase.err.Code, hook.LastEntry().Data["errorCode"], "Error in test case %v", n)
		assert.Equal(t, testcase.httpMethod, hook.LastEntry().Data["httpMethod"], "Error in test case %v", n)
		assert.Equal(t, testcase.httpAddress+testcase.httpURI, hook.LastEntry().Data["httpURI"], "Error in test case %v", n)
		assert.Equal(t, "", hook.LastEntry().Data["httpRemoteAddress"], "Error in test case %v", n)
		assert.Equal(t, testcase.httpStatusCode, hook.LastEntry().Data["httpStatusCode"], "Error in test case %v", n)
		hook.Reset()
	}
}
