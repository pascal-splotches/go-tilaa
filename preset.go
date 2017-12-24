package go_tilaa

const presetsBasePath = "presets"

type PresetServiceInterface interface {
	List() (*Presets, error)
}

type PresetService struct {
	client *Client
}

var _ PresetServiceInterface = &PresetService{}

type Presets struct {
	Ram struct {
		Sizes []int `json:"sizes"`
	} `json:"ram"`

	Storage []struct {
		Type  string `json:"type"`
		Sizes []int  `json:"sizes"`
	} `json:"storage"`

	client *Client `json:"-"`
}

type PresetsResponse struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message,omitempty"`
	Presets Presets        `json:"presets"`
}

func (service *PresetService) List() (*Presets, error) {
	// TODO: Runtime cache presets
	var response PresetsResponse

	_, err := service.client.Get(presetsBasePath, &response)

	if err != nil {
		return nil, err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	presets := response.Presets

	presets.client = service.client

	return &presets, err
}

// TODO: Add validation helpers to Presets structure

func NewPresets() *Presets {
	return &Presets{}
}