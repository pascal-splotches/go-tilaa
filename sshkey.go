package go_tilaa

import (
	"time"
	"net/url"
	"strconv"
	"fmt"
)

const sshKeyBasePath = "ssh_keys"

type SshKeyServiceInterface interface {
	List() (*[]SshKey, error)
	Add(*SshKey) (*SshKey, error)
	View(int) (*SshKey, error)
	Edit(*SshKey) (*SshKey, error)
	Delete(*SshKey) error
}

type SshKeyService struct {
	client *Client
}

var _ SshKeyServiceInterface = &SshKeyService{}

type SshKey struct {
	Id       int       `json:"id"`
	UserId   int       `json:"user_id"`
	Label    string    `json:"label"`
	Key      string    `json:"key"`
	Created  time.Time `json:"created,string"`
	Modified time.Time `json:"modified,string"`

	client *Client `json:"-"`
}

type SshKeyResponse struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message,omitempty"`
	SshKey  SshKey         `json:"ssh_key"`
}

type SshKeysResponse struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message,omitempty"`
	SshKeys []SshKey       `json:"ssh_keys"`
}

func (service *SshKeyService) List() (*[]SshKey, error) {
	var response SshKeysResponse

	_, err := service.client.Get(sshKeyBasePath, &response)

	if err != nil {
		return nil, err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	sshKeys := response.SshKeys

	for i := range sshKeys {
		sshKeys[i].client = service.client
	}

	return &sshKeys, err
}

func (service *SshKeyService) Add(sshKey *SshKey) (*SshKey, error) {
	if err := sshKey.Validate(); err != nil {
		return NewSshKey(service.client), err
	}

	payload := sshKey.Payload()

	var response StatusResponse

	_, err := service.client.Post(sshKeyBasePath, payload, &response)

	if err != nil {
		return NewSshKey(service.client), err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	// TODO: SSH Key endpoint does not return ID of the created key

	return sshKey, err
}

func (service *SshKeyService) View(sshKeyId int) (*SshKey, error) {
	var response SshKeyResponse

	_, err := service.client.Get(service.path(strconv.Itoa(sshKeyId)), &response)

	if err != nil {
		return NewSshKey(service.client), err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	sshKey := response.SshKey

	sshKey.client = service.client

	return &sshKey, err
}

func (service *SshKeyService) Edit(sshKey *SshKey) (*SshKey, error) {
	if err := sshKey.Validate(); err != nil {
		return sshKey, err
	}

	payload := sshKey.Payload()

	var response StatusResponse

	_, err := service.client.Post(service.path(strconv.Itoa(sshKey.Id)), payload, &response)

	if err != nil {
		return sshKey, err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	return sshKey, err
}

func (service *SshKeyService) Delete(sshKey *SshKey) error {
	var response StatusResponse

	_, err := service.client.Delete(service.path(strconv.Itoa(sshKey.Id)), &response)

	if err != nil {
		return err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	return err
}

func (service *SshKeyService) path(path string) string {
	return fmt.Sprintf("%s/%s", sshKeyBasePath, path)
}

func (sshKey *SshKey) Create() error {
	_, err := sshKey.client.SshKey.Add(sshKey)

	return err
}

func (sshKey *SshKey) Commit() error {
	if sshKey.Id == 0 {
		return NewSshKeyNotCreatedError()
	}

	_, err := sshKey.client.SshKey.Edit(sshKey)

	return err
}

func (sshKey *SshKey) Delete() error {
	if sshKey.Id == 0 {
		return NewSshKeyNotCreatedError()
	}

	err := sshKey.client.SshKey.Delete(sshKey)

	return err
}

func (sshKey *SshKey) Payload() *url.Values {
	return &url.Values {
		"user_id": {strconv.Itoa(sshKey.UserId)},
		"label":   {sshKey.Label},
		"key":     {sshKey.Key},
	}
}

func (sshKey *SshKey) Validate() error {
	// TODO: Implement validation for all the fields
	return nil
}

func NewSshKey(client *Client) *SshKey {
	return &SshKey{client: client}
}