package go_tilaa

const templatesBasePath = "templates"

type TemplateServiceInterface interface {
	List() (*[]Template, error)
}

type TemplateService struct {
	client *Client
}

var _ TemplateServiceInterface = &TemplateService{}

// TODO: Figure out why some templates return 0 storage and/or 0 ram. Does this mean the preset minimum is the minimum?
// TODO: Figure out why template 166 has no name
type Template struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Ram     int    `json:"ram"`
	Storage int    `json:"storage"`

	client *Client `json:"-"`
}

type TemplatesResponse struct {
	Status    ResponseStatus `json:"status"`
	Message   string         `json:"message,omitempty"`
	Templates []Template     `json:"templates"`
}

func (service *TemplateService) List() (*[]Template, error) {
	var response TemplatesResponse

	_, err := service.client.Get(templatesBasePath, &response)

	if err != nil {
		return nil, err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	templates := response.Templates

	for i := range templates {
		templates[i].client = service.client
	}

	return &templates, err
}
