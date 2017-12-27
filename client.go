package go_tilaa

import (
	"net/http"
	"net/url"
	"encoding/json"
	"strings"
	"io"
	"fmt"
)

const (
	ApiVersion = "v1"

	BaseUrl = "https://api.tilaa.com"

	Accept = "application/json"

	UserAgent = "go-tilaa"

	ContentTypeFormUrlEncoded = "application/x-www-form-urlencoded"
)

var (
	headerUserAgentKey = http.CanonicalHeaderKey("User-Agent")
	acceptKey          = http.CanonicalHeaderKey("Accept")
	contentTypeKey     = http.CanonicalHeaderKey("Content-Type")
)

type Client struct {
	BaseUrl     *url.URL
	UserAgent   string
	Credentials BasicAuth

	httpClient *http.Client

	VirtualMachine VirtualMachineServiceInterface
	Snapshot       SnapshotServiceInterface
	Preset         PresetServiceInterface
	Template       TemplateServiceInterface
	Site           SiteServiceInterface
	Metadata       MetadataServiceInterface
	SshKey         SshKeyServiceInterface
}

type BasicAuth struct {
	UserName string
	Password string
}

func New(username string, password string) *Client {
	return createClient(username, password)
}

func (client *Client) SetBasicAuth(username string, password string) {
	client.Credentials.UserName = username
	client.Credentials.Password = password
}

func (client *Client) Get(path string, result interface{}) (*http.Response, error) {
	request, err := client.newRequest(http.MethodGet, path, nil, "")

	if err != nil {
		return nil, err
	}

	return client.do(request, result)
}

func (client *Client) Post(path string, formData *url.Values, result interface{}) (*http.Response, error) {
	request, err := client.newRequest(http.MethodPost, path, strings.NewReader(formData.Encode()), ContentTypeFormUrlEncoded)

	if err != nil {
		return nil, err
	}

	return client.do(request, result)
}

func (client *Client) Delete(path string, result interface{}) (*http.Response, error) {
	request, err := client.newRequest(http.MethodDelete, path, nil, "")

	if err != nil {
		return nil, err
	}

	return client.do(request, result)
}

func (client *Client) newRequest(method string, path string, body io.Reader, contentType string) (*http.Request, error) {
	path = fmt.Sprintf("%s/%s", ApiVersion, path)

	relativePath := &url.URL{Path: path}
	requestUrl   := client.BaseUrl.ResolveReference(relativePath)

	request, err := http.NewRequest(method, requestUrl.String(), body)

	if err != nil {
		return nil, NewApiRequestError(err.Error())
	}

	request.Header.Set(headerUserAgentKey, client.UserAgent)
	request.Header.Set(acceptKey, Accept)

	if body != nil && contentType != "" {
		request.Header.Set(contentTypeKey, contentType)
	}

	request.SetBasicAuth(client.Credentials.UserName, client.Credentials.Password)

	return request, nil
}

func (client *Client) do(request *http.Request, result interface{}) (*http.Response, error) {
	response, err := client.httpClient.Do(request)

	if err != nil {
		return nil, NewApiRequestError(err.Error())
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return nil, NewInvalidCredentialsError(client.Credentials)
	}

	if response.StatusCode != http.StatusOK {
		return nil, NewApiRequestError(fmt.Sprintf("[%s] Request returned with non-200 status", response.Status))
	}

	err = json.NewDecoder(response.Body).Decode(result)

	if err != nil {
		return nil, NewResultsDecoderError(err, response)
	}

	return response, err
}

type ResponseStatus string

const (
	ResponseOk    = "OK"
	ResponseError = "ERROR"
)

type StatusResponse struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message,omitempty"`
}