package api

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// Log is API global logger for all domains
var Log *logrus.Logger

// LogOperation logs an action with request identifier and user
func LogOperation(requestID string, userID string, message string) {
	fields := getLogFields(requestID, userID, "", nil, 0, nil)
	Log.WithFields(fields).Info(message)
}

// LogOperation logs a warning action with request identifier and user
func LogOperationWarn(requestID string, userID string, message string) {
	fields := getLogFields(requestID, userID, "", nil, 0, nil)
	Log.WithFields(fields).Warn(message)
}

// LogErrorMessage logs an error action with request identifier and user
func LogOperationError(requestID string, userID string, err *Error) {
	fields := getLogFields(requestID, userID, "", nil, 0, err)
	Log.WithFields(fields).Error(err.Message)
}

// TransactionRequestLog logs a request transaction received with http request, user and request identifier
func TransactionRequestLog(requestID string, userID string, r *http.Request) {
	fields := getLogFields(requestID, userID, "", r, 0, nil)
	Log.WithFields(fields).Info("")
}

// TransactionResponseErrorLog logs a response error transaction with http request, user, request identifier and status code
func TransactionResponseErrorLog(requestID string, userID string, r *http.Request, status int, err *Error) {
	fields := getLogFields(requestID, userID, "", r, status, err)
	Log.WithFields(fields).Error(err.Message)
}

// TransactionProxyLog logs a request transaction received with user, worker request identifier and request identifier
func TransactionProxyLog(requestID string, workerRequestID string, r *http.Request, msg string) {
	fields := getLogFields(requestID, "", workerRequestID, r, 0, nil)
	Log.WithFields(fields).Info(msg)
}

// TransactionProxyErrorLog logs an error received with user, worker request identifier, proxy request identifier and status code
func TransactionProxyErrorLogWithStatus(requestID string, workerRequestID string, r *http.Request, status int, err *Error) {
	fields := getLogFields(requestID, "", workerRequestID, r, status, err)
	Log.WithFields(fields).Error(err.Message)
}

// getLogFields returns a map with fields according to input parameters
func getLogFields(
	requestID string, userID string, workerRequestID string, // Basic Info
	r *http.Request, status int, // HTTP info
	err *Error,
) map[string]interface{} {

	fields := make(map[string]interface{})
	// Create fields according input parameters
	if requestID != "" {
		fields["requestID"] = requestID
	}
	if userID != "" {
		fields["user"] = userID
	}
	if workerRequestID != "" {
		fields["workerRequestID"] = workerRequestID
	}
	if r != nil {
		// TODO: X-Forwarded headers
		fields["httpMethod"] = r.Method
		fields["httpURI"] = r.URL.EscapedPath()
		fields["httpRemoteAddress"] = r.RemoteAddr
	}
	if status != 0 {
		fields["httpStatusCode"] = status
	}
	if err != nil {
		fields["errorCode"] = err.Code
	}

	return fields
}
