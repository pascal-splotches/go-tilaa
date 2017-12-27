package go_tilaa

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const metadataBasePath = "metadata"

type MetadataServiceInterface interface {
	List() (*[]Metadata, error)
	Add(*Metadata) (*Metadata, error)
	View(int) (*Metadata, error)
	Edit(*Metadata) (*Metadata, error)
	Delete(*Metadata) error
}

type MetadataService struct {
	client *Client
}

var _ MetadataServiceInterface = &MetadataService{}

type Metadata struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	UserData string    `json:"user_data"`
	Created  time.Time `json:"created,string"`
	Modified time.Time `json:"modified,string"`

	client *Client `json:"-"`
}

type MetadataResponse struct {
	Status   ResponseStatus `json:"status"`
	Message  string         `json:"message,omitempty"`
	Metadata Metadata       `json:"metadata"`
}

type MetadatasResponse struct {
	Status   ResponseStatus `json:"status"`
	Message  string         `json:"message,omitempty"`
	Metadata []Metadata     `json:"metadata"`
}

type NewMetadataResponse struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message,omitempty"`
	Id      int            `json:"id"`
}

func (service *MetadataService) List() (*[]Metadata, error) {
	var response MetadatasResponse

	_, err := service.client.Get(metadataBasePath, &response)

	if err != nil {
		return nil, err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	metadata := response.Metadata

	for i := range metadata {
		metadata[i].client = service.client
	}

	return &metadata, err
}

func (service *MetadataService) Add(metadata *Metadata) (*Metadata, error) {
	if err := metadata.Validate(); err != nil {
		return NewMetadata(service.client), err
	}

	payload := metadata.Payload()

	var response NewMetadataResponse

	_, err := service.client.Post(metadataBasePath, payload, &response)

	if err != nil {
		return NewMetadata(service.client), err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	metadata.Id = response.Id

	return metadata, err
}

func (service *MetadataService) View(metadataId int) (*Metadata, error) {
	var response MetadataResponse

	_, err := service.client.Get(service.path(strconv.Itoa(metadataId)), &response)

	if err != nil {
		return NewMetadata(service.client), err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	metadata := response.Metadata

	metadata.client = service.client

	return &metadata, err
}

func (service *MetadataService) Edit(metadata *Metadata) (*Metadata, error) {
	if err := metadata.Validate(); err != nil {
		return metadata, err
	}

	payload := metadata.Payload()

	var response StatusResponse

	_, err := service.client.Post(service.path(strconv.Itoa(metadata.Id)), payload, &response)

	if err != nil {
		return metadata, err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	return metadata, err
}

func (service *MetadataService) Delete(metadata *Metadata) error {
	var response StatusResponse

	_, err := service.client.Delete(service.path(strconv.Itoa(metadata.Id)), &response)

	if err != nil {
		return err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	return err
}

func (service *MetadataService) path(path string) string {
	return fmt.Sprintf("%s/%s", metadataBasePath, path)
}

func (metadata *Metadata) Create() error {
	_, err := metadata.client.Metadata.Add(metadata)

	return err
}

func (metadata *Metadata) Commit() error {
	if metadata.Id == 0 {
		return NewMetadataNotCreatedError()
	}

	_, err := metadata.client.Metadata.Edit(metadata)

	return err
}

func (metadata *Metadata) Delete() error {
	if metadata.Id == 0 {
		return NewMetadataNotCreatedError()
	}

	err := metadata.client.Metadata.Delete(metadata)

	return err
}

func (metadata *Metadata) Payload() *url.Values {
	return &url.Values{
		"name":      {metadata.Name},
		"user_data": {metadata.UserData},
	}
}

func (metadata *Metadata) Validate() error {
	// TODO: Implement validation for all the fields
	return nil
}

func NewMetadata(client *Client) *Metadata {
	return &Metadata{client: client}
}
