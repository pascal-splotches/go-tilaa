package go_tilaa

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"net/url"
	"reflect"
)

var (
	basicAuth = BasicAuth{
		UserName: "test",
		Password: "test123",
	}
)

func TestClient(t *testing.T) {
	testServer := createApiServer(t)
	defer testServer.Close()

	testClient := New(basicAuth.UserName, basicAuth.Password)

	testClient.BaseUrl, _ = url.Parse(testServer.URL)

	var response StatusResponse
	rawResponse, err := testClient.Get("", &response)

	assertError(t, err)
	assertEqual(t, http.StatusOK, rawResponse.StatusCode)
	assertNotNil(t, response)
	assertEqual(t, ResponseOk, response.Status)
	assertEqual(t, "Test Response", response.Message)
}

func TestClient_InvalidAuth(t *testing.T) {
	testServer := createApiServer(t)
	defer testServer.Close()

	testClient := New("invalid", "invalid")

	testClient.BaseUrl, _ = url.Parse(testServer.URL)

	var response StatusResponse
	rawResponse, err := testClient.Get("", &response)

	assertNotNil(t, err)

	if _, ok := err.(*InvalidCredentialsError); !ok {
		assertInvalidType(t, "InvalidCredentialsError", reflect.TypeOf(err))
	}

	assertNil(t, rawResponse)
}

func TestClient_SetBasicAuth(t *testing.T) {
	testServer := createApiServer(t)
	defer testServer.Close()

	testClient := New("invalid", "invalid")

	testClient.BaseUrl, _ = url.Parse(testServer.URL)

	var response StatusResponse
	_, err := testClient.Get("", &response)

	assertNotNil(t, err)

	testClient.SetBasicAuth(basicAuth.UserName, basicAuth.Password)

	_, err = testClient.Get("", &response)

	assertNil(t, err)
}

func TestClient_Get(t *testing.T) {
	testServer := createApiServer(t)
	defer testServer.Close()

	testClient := New(basicAuth.UserName, basicAuth.Password)

	testClient.BaseUrl, _ = url.Parse(testServer.URL)

	var response StatusResponse
	rawResponse, err := testClient.Get("", &response)

	assertError(t, err)
	assertEqual(t, http.StatusOK, rawResponse.StatusCode)
	assertNotNil(t, response)
	assertEqual(t, ResponseOk, response.Status)
}

func TestClient_Get2(t *testing.T) {
	testServer := createApiServer(t)
	defer testServer.Close()

	testClient := New(basicAuth.UserName, basicAuth.Password)

	testClient.BaseUrl, _ = url.Parse(testServer.URL)

	var response StatusResponse
	rawResponse, err := testClient.Get("response-error", &response)

	assertNotNil(t, err)

	if _, ok := err.(*ApiRequestError); !ok {
		assertInvalidType(t, "ApiRequestError", reflect.TypeOf(err))
	}

	assertNil(t, rawResponse)
	assertEqual(t, ResponseStatus(""), response.Status)
}

func TestClient_Get3(t *testing.T) {
	testServer := createApiServer(t)
	defer testServer.Close()

	testClient := New(basicAuth.UserName, basicAuth.Password)

	testClient.BaseUrl, _ = url.Parse(testServer.URL)

	var response StatusResponse
	rawResponse, err := testClient.Get("invalid-response", &response)

	assertNotNil(t, err)

	if _, ok := err.(*ResultsDecoderError); !ok {
		assertInvalidType(t, "ResultsDecoderError", reflect.TypeOf(err))
	}

	assertNil(t, rawResponse)
	assertEqual(t, ResponseStatus(""), response.Status)
}

func createApiServer(t *testing.T) *httptest.Server {
	testServer := createTestServer(func(writer http.ResponseWriter, request *http.Request) {
		t.Logf("Method: %v", request.Method)
		t.Logf("Path: %v", request.URL.Path)

		if !handleBasicAuth(t, writer, request) {
			writer.WriteHeader(http.StatusUnauthorized)
			_, _ = writer.Write([]byte("{\"status\":\"ERROR\",\"message\":\"Unauthorized\"}"))
			return
		}

		switch request.Method {
		case http.MethodGet:
			switch request.URL.Path {
			case "/v1/":
				_, _ = writer.Write([]byte("{\"status\":\"OK\",\"message\":\"Test Response\"}"))
			case "/v1/response-error":
				writer.WriteHeader(http.StatusNotFound)
				_, _ = writer.Write([]byte("{\"status\":\"ERROR\",\"message\":\"Test Non-200 Status\"}"))
			case "/v1/invalid-response":
				_, _ = writer.Write([]byte("{\"status\":\"OK\",}"))
			}
		}
	})

	return testServer
}

func handleBasicAuth(t *testing.T, writer http.ResponseWriter, request *http.Request) bool {
	username, password, _ := request.BasicAuth()

	if basicAuth.UserName == username && basicAuth.Password == password {
		return true
	}

	return false
}

func createTestServer(function func(writer http.ResponseWriter, request *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(function))
}

func assertNil(t *testing.T, v interface{}) {
	if !isNil(v) {
		t.Errorf("[%v] was expected to be nil", v)
	}
}

func assertNotNil(t *testing.T, v interface{}) {
	if isNil(v) {
		t.Errorf("[%v] was expected to be non-nil", v)
	}
}

func assertError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Error occurred (%s) [%v]", err.Error(), err)
	}
}

func assertEqual(t *testing.T, expected, got interface{}) (r bool) {
	if !equal(expected, got) {
		t.Errorf("Expected [%#v], got [%#v]", expected, got)
	}

	return
}

func assertNotEqual(t *testing.T, expected, got interface{}) (r bool) {
	if equal(expected, got) {
		t.Errorf("expected [%v], got [%v]", expected, got)
	}

	return
}

func assertType(t *testing.T, expected, got interface{}) {
	if reflect.TypeOf(expected) != reflect.TypeOf(got) {
		assertInvalidType(t, reflect.TypeOf(expected), reflect.TypeOf(got))
	}
}

func assertInvalidType(t *testing.T, expected, got interface{}) {
	t.Errorf("[%v] was expected to be of type [%v]", got, expected)
}

func isNil(v interface{}) bool {
	if v == nil {
		return true
	}

	rv   := reflect.ValueOf(v)
	kind := rv.Kind()

	if kind >= reflect.Chan && kind <= reflect.Slice && rv.IsNil() {
		return true
	}

	return false
}

func equal(expected, got interface{}) bool {
	return reflect.DeepEqual(expected, got)
}