package go_tilaa

import (
	"net/http"
	"net/url"
)

var DefaultClient *Client

func New(username string, password string) *Client {
	return createClient(username, password)
}

func SetBasicAuth(username string, password string) {
	DefaultClient.SetBasicAuth(username, password)
}

func createClient(username string, password string) *Client {
	httpClient := &http.Client{}

	client := &Client{
		httpClient: httpClient,
	}

	client.SetBasicAuth(username, password)
	client.BaseUrl, _ = url.Parse(BaseUrl)
	client.UserAgent  = UserAgent + "/" + ApiVersion

	client.VirtualMachine = &VirtualMachineService{client: client}
	client.Snapshot       = &SnapshotService{client: client}
	client.Preset         = &PresetService{client: client}
	client.Template       = &TemplateService{client: client}
	client.Site           = &SiteService{client: client}
	client.Metadata       = &MetadataService{client: client}
	client.SshKey         = &SshKeyService{client: client}

	return client
}

func init() {
	DefaultClient = New("", "")
}