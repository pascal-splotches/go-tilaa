package go_tilaa

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const snapshotBasePath = "snapshots"

type SnapshotServiceInterface interface {
	List() (*[]Snapshot, error)
	Add(*VirtualMachine, string, bool, bool) (*Snapshot, error)
	View(int) (*Snapshot, error)
	Rename(*Snapshot, string) (*Snapshot, error)
	Delete(*Snapshot) error
	Restore(*VirtualMachine, *Snapshot) (*VirtualMachine, error)
}

type SnapshotService struct {
	client *Client
}

var _ SnapshotServiceInterface = &SnapshotService{}

type Snapshot struct {
	Id       int            `json:"id"`
	Name     string         `json:"name"`
	Storage  int            `json:"storage"`
	Ram      int            `json:"ram"`
	Template Template       `json:"template"`
	Status   SnapshotStatus `json:"status"`
	Created  time.Time      `json:"created,string"`

	client *Client `json:"-"`
}

type SnapshotStatus string

// TODO: Figure out snapshot statuses
const (
	SnapshotStatusSuccess = "success"
)

type SnapshotResponse struct {
	Status   ResponseStatus `json:"status"`
	Message  string         `json:"message,omitempty"`
	Snapshot Snapshot       `json:"snapshot"`
}

type SnapshotsResponse struct {
	Status    ResponseStatus `json:"status"`
	Message   string         `json:"message,omitempty"`
	Snapshots []Snapshot     `json:"snapshots"`
}

func (service *SnapshotService) List() (*[]Snapshot, error) {
	var response SnapshotsResponse

	_, err := service.client.Get(snapshotBasePath, &response)

	if err != nil {
		return nil, err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	snapshots := response.Snapshots

	for i := range snapshots {
		snapshots[i].client = service.client
	}

	return &snapshots, err
}

func (service *SnapshotService) Add(machine *VirtualMachine, name string, online bool, overwrite bool) (*Snapshot, error) {
	return service.client.VirtualMachine.CreateSnapshot(machine, name, online, overwrite)
}

func (service *SnapshotService) View(snapshotId int) (*Snapshot, error) {
	var response SnapshotResponse

	_, err := service.client.Get(service.path(strconv.Itoa(snapshotId)), &response)

	if err != nil {
		return NewSnapshot(service.client), err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	snapshot := response.Snapshot

	snapshot.client = service.client

	return &snapshot, err
}

func (service *SnapshotService) Rename(snapshot *Snapshot, name string) (*Snapshot, error) {
	payload := &url.Values{
		"name": {name},
	}

	var response StatusResponse

	_, err := service.client.Post(service.path(strconv.Itoa(snapshot.Id)), payload, &response)

	if err != nil {
		return snapshot, err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	snapshot.Name = name

	return snapshot, err
}

func (service *SnapshotService) Delete(snapshot *Snapshot) error {
	var response StatusResponse

	_, err := service.client.Delete(service.path(strconv.Itoa(snapshot.Id)), &response)

	if err != nil {
		return err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	return err
}

func (service *SnapshotService) Restore(machine *VirtualMachine, snapshot *Snapshot) (*VirtualMachine, error) {
	return service.client.VirtualMachine.RestoreSnapshot(machine, snapshot)
}

func (service *SnapshotService) path(path string) string {
	return fmt.Sprintf("%s/%s", snapshotBasePath, path)
}

func (snapshot *Snapshot) Create(machine *VirtualMachine, online bool, overwrite bool) error {
	_, err := machine.CreateSnapshot(snapshot.Name, online, overwrite)

	// TODO: Endpoint does not return Snapshot ID after creation

	return err
}

func (snapshot *Snapshot) Rename(name string) error {
	if snapshot.Id == 0 {
		return NewSnapshotNotCreatedError()
	}

	_, err := snapshot.client.Snapshot.Rename(snapshot, name)

	if err == nil {
		snapshot.Name = name
	}

	return err
}

func (snapshot *Snapshot) Delete() error {
	if snapshot.Id == 0 {
		return NewSnapshotNotCreatedError()
	}

	err := snapshot.client.Snapshot.Delete(snapshot)

	return err
}

func (snapshot *Snapshot) RestoreToVirtualMachine(machine *VirtualMachine) error {
	if snapshot.Id == 0 {
		return NewSnapshotNotCreatedError()
	}

	if machine.Id == 0 {
		return NewVirtualMachineNotCreatedError()
	}

	_, err := machine.client.VirtualMachine.RestoreSnapshot(machine, snapshot)

	return err
}

func NewSnapshot(client *Client) *Snapshot {
	return &Snapshot{client: client}
}
