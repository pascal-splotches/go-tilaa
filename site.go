package go_tilaa

const sitesBasePath = "sites"

type SiteServiceInterface interface {
	List() (*[]Site, error)
}

type SiteService struct {
	client *Client
}

var _ SiteServiceInterface = &SiteService{}

type Site struct {
	Id   int    `json:"id"`
	Name string `json:"name"`

	client *Client `json:"-"`
}

type SitesResponse struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message,omitempty"`
	Sites   []Site         `json:"sites"`
}

func (service *SiteService) List() (*[]Site, error) {
	var response SitesResponse

	_, err := service.client.Get(sitesBasePath, &response)

	if err != nil {
		return nil, err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	sites := response.Sites

	for i := range sites {
		sites[i].client = service.client
	}

	return &sites, err
}