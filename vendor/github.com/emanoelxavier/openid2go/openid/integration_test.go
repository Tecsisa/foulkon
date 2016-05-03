// +build integration

package openid_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/emanoelxavier/openid2go/openid"
)

const authenticatedMessage string = "Congrats, you are authenticated!"
const authenticatedMessageWithUser string = "Congrats, you are authenticated by the provider %v!"

// The idToken flag must have a valid ID Token issued by any OIDC provider.
var idToken = flag.String("idToken", "", "a valid id token")

// The issuer and cliendID  flags must be valid.
var issuer = flag.String("issuer", "", "the OP issuer")
var clientID = flag.String("clientID", "", "the client ID registered with the OP")

var mux *http.ServeMux

// The authenticateHandler is registered behind the openid.Authenticate middleware.
func authenticatedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, authenticatedMessage)
}

// The authenticateHandlerWithUser is registered behind the openid.AuthenticateUser middleware.
func authenticatedHandlerWithUser(u *openid.User, w http.ResponseWriter, r *http.Request) {
	iss := u.Issuer

	// Workaround for Google OP since it has a bug causing the 'iss' claim to miss the 'https://'
	if strings.HasPrefix(*issuer, "https://") && !strings.HasPrefix(iss, "https://") {
		iss = "https://" + iss
	}

	fmt.Fprintf(w, "%v User: %+v.\n", fmt.Sprintf(authenticatedMessageWithUser, iss), u)
}

// The init func initializes the openid.Configuration and the server routes.
func init() {
	mux = http.NewServeMux()
	config, err := openid.NewConfiguration(openid.ProvidersGetter(getProviders))

	if err != nil {
		fmt.Println("Error whe creating the configuration for the openid middleware.", err)
	}

	mux.Handle("/authn", openid.Authenticate(config, http.HandlerFunc(authenticatedHandler)))
	mux.Handle("/user", openid.AuthenticateUser(config, openid.UserHandlerFunc(authenticatedHandlerWithUser)))
}

// Validates that a valid ID Token results in a successful authentication of the user.
func Test_Authenticate_ValidIDToken(t *testing.T) {
	// Arrange.
	server := httptest.NewServer(mux)
	defer server.Close()

	// Act.
	res, code, err := executeRequest(server.URL, "/authn", *idToken)

	// Assert.
	validateResponse(t, err, res, code, authenticatedMessage, http.StatusOK)
}

// Validates that an ID Token signed by an unknown key identifier results in HTTP Status Unauthorized.
func Test_Authenticate_InvalidIDTokenKeyID(t *testing.T) {
	// Arrange.
	server := httptest.NewServer(mux)
	defer server.Close()
	et := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImNiMzVkMTZjZmI4MWY2ZTUzZDk5YTBmODg4YjRhZTgyNWE3MWU1Y2MifQ.eyJpc3MiOiJhY2NvdW50cy5nb29nbGUuY29tIiwiYXRfaGFzaCI6Im1iVVZpRlFReUFPX2Y1YlR0alVvREEiLCJhdWQiOiI0MDc0MDg3MTgxOTIuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMTI2Nzg1OTg3MTA3MDYxNzA2NDkiLCJhenAiOiI0MDc0MDg3MTgxOTIuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJpYXQiOjE0NTA5MjkxMjAsImV4cCI6MTQ1MDkzMjcyMH0.f5toakDvtU3Tqt71uDgIACrac8mGM4K8HQ1Fyw9jaUdxonEu_Bww-UNKjPD6tKAe7AzVJzfKOzzcvJygMRfQ72u4wsljhQV3i6-cJmpMj4S5HQoleV4GqNHq-84KNEvFv_4IT7wIEdu0kEpRygt9lhysvFXxGfkR6TpTr50W8yo4T0EfRVXafXhNMX5uNkVJ2PnHQXYvwiZP2_hJ6KTw9RN_4mDJIOHJBAXRXnBe9hJ7pDfED1v_ayOJGmq5PGeAO34Rz7FDf4Awf8DOoMiRVSi3SwE_pRHwBRp2WsrvKpSdUToXyOmBEzPYubEJIcYPiJR4uXgPEOynV0i993NGQA"

	// Changing the client ID and issuer to match the ones in the token.
	cID := "407408718192.apps.googleusercontent.com"
	cIDp := &cID
	cIDp, clientID = clientID, cIDp
	iss := "accounts.google.com"
	issp := &iss
	issp, issuer = issuer, issp

	// Act.
	res, code, err := executeRequest(server.URL, "/user", et)

	// Assert.
	validateResponse(t, err, res, code, "", http.StatusUnauthorized)

	// Clean up.
	// Swapping back the values for the next test.
	clientID = cIDp
	issuer = issp
}

// Validates that an ID Token issued for an audience that is not registered as an
// allowed client ID returns unauthorized. This test also demonstrates that configuration
// changes takes effect without the need of a service start.
func Test_Authenticate_ChangeClientID_InvalidateIDToken(t *testing.T) {
	// Arrange.
	server := httptest.NewServer(mux)
	defer server.Close()

	// Changing client ID for another value to invalidate the token.
	cID := "newClientID"
	cIDp := &cID
	cIDp, clientID = clientID, cIDp

	// Act.
	res, code, err := executeRequest(server.URL, "/user", *idToken)

	// Assert.
	validateResponse(t, err, res, code, "", http.StatusUnauthorized)

	// Clean up.
	// Swapping back the values for the next test.
	clientID = cIDp
}

// Validates that a valid ID Token results in a successful authentication of the user.
// And that the user information extracted from the token is made available to the rest of the
// application stack.
func Test_AuthenticateUser_ValidIDToken(t *testing.T) {
	// Arrange.
	server := httptest.NewServer(mux)
	defer server.Close()

	// Act.
	res, code, err := executeRequest(server.URL, "/user", *idToken)

	// Assert.
	validateResponse(t, err, res, code, fmt.Sprintf(authenticatedMessageWithUser, *issuer), http.StatusOK)
}

// getProviders is used as the GetProvidersFunc when creating a new openid.Configuration.
func getProviders() ([]openid.Provider, error) {
	return []openid.Provider{{Issuer: *issuer, ClientIDs: []string{*clientID}}}, nil
}

func executeRequest(u string, r string, t string) (string, int, error) {
	var res string
	var code int
	client := http.DefaultClient

	req, err := http.NewRequest("GET", u+r, nil)
	if err != nil {
		return res, code, err
	}

	req.Header.Add("Authorization", "Bearer "+t)

	resp, err := client.Do(req)
	if err != nil {
		return res, code, err
	}

	msg, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return res, code, err
	}

	res = string(msg)
	code = resp.StatusCode

	return res, code, nil
}

func validateResponse(t *testing.T, e error, r string, c int, er string, ec int) {
	if e != nil {
		t.Error(e)
	}

	if er != "" && !strings.HasPrefix(r, er) {
		t.Error("Expected response with prefix:", er, "but got:", r)
	} else {
		t.Log(r)
	}

	if c != ec {
		t.Error("Expected HTTP status", ec, "but got", c)
	}
}
