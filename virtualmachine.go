package go_tilaa

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"
)

const virtualMachinesBasePath = "virtual_machines"

type VirtualMachineServiceInterface interface {
	List() (*[]VirtualMachine, error)
	Add(*VirtualMachine) (*VirtualMachine, error)
	AddFromSnapshot(*VirtualMachine, *Snapshot) (*VirtualMachine, error)
	View(int) (*VirtualMachine, error)
	Edit(*VirtualMachine) (*VirtualMachine, error)
	Cancel(*VirtualMachine, *time.Time) (*VirtualMachine, error)
	UndoCancellation(*VirtualMachine) error
	GetCancelDates(*VirtualMachine) (*[]time.Time, error)
	RunTask(string, *VirtualMachine) (*VirtualMachine, error)
	Reinstall(*VirtualMachine) (*VirtualMachine, error)
	CreateSnapshot(*VirtualMachine, string, bool, bool) (*Snapshot, error)
	RestoreSnapshot(*VirtualMachine, *Snapshot) (*VirtualMachine, error)
}

type VirtualMachineService struct {
	client *Client
}

var _ VirtualMachineServiceInterface = &VirtualMachineService{}

type VirtualMachine struct {
	Id        int                  `json:"id"`
	Name      string               `json:"name"`
	Cpu       Cpu                  `json:"cpu"`
	Ram       int                  `json:"ram"`
	Storage   Storage              `json:"storage"`
	Site      Site                 `json:"site"`
	Template  Template             `json:"template"`
	Network   []Network            `json:"network"`
	Admin     Admin                `json:"admin"`
	Managed   bool                 `json:"is_managed"`
	Locked    bool                 `json:"locked"`
	Status    VirtualMachineStatus `json:"status"`
	Created   *time.Time           `json:"created,string"`
	Cancelled *time.Time           `json:"cancelled,string"`

	client   *Client `json:"-"`
	original struct {
		Name     string
		Ram      int
		Storage  int
		Template int
		Cpu      Cpu
	} `json:"-"`
}

const (
	TaskStart    = "start"
	TaskStop     = "stop"
	TaskRestart  = "restart"
	TaskPowerOff = "poweroff"
	TaskRescue   = "rescue"
)

type Cpu struct {
	Cores int `json:"count"`
	Cap   int `json:"cap"`
}

type StorageType string

const (
	StorageTypeHdd StorageType = "hdd"
	StorageTypeSsd StorageType = "ssd"
)

type Storage struct {
	Size int         `json:"size"`
	Type StorageType `json:"type"`
}

type NetworkFamily int

const (
	NetworkFamilyIpv4 NetworkFamily = 4
	NetworkFamilyIpv6 NetworkFamily = 6
)

type Network struct {
	Id      int           `json:"id"`
	Family  NetworkFamily `json:"family,int"`
	Address net.IP        `json:"address,string"`
	DnsName string        `json:"dns_name"`
}

type Admin struct {
	Account         string `json:"account"`
	InitialPassword string `json:"initial_password"`
}

type VirtualMachineStatus string

const (
	VirtualMachineStatusPending                 = "pending"
	VirtualMachineStatusCreateFailed            = "create_failed"
	VirtualMachineStatusDestroyed               = "destroyed"
	VirtualMachineStatusCreating                = "creating"
	VirtualMachineStatusStopped                 = "stopped"
	VirtualMachineStatusStarting                = "starting"
	VirtualMachineStatusRestarting              = "restarting"
	VirtualMachineStatusRunning                 = "running"
	VirtualMachineStatusRunning_Rescue          = "running_rescue"
	VirtualMachineStatusStopping                = "stopping"
	VirtualMachineStatusResize_Failed           = "resize_failed"
	VirtualMachineStatusDestroying              = "destroying"
	VirtualMachineStatusDestroy_Failed          = "destroy_failed"
	VirtualMachineStatusResizing                = "resizing"
	VirtualMachineStatusMigrating               = "migrating"
	VirtualMachineStatusMigrating_Live          = "live_migrating"
	VirtualMachineStatusMigrate_Failed          = "migrate_failed"
	VirtualMachineStatusPaused                  = "paused"
	VirtualMachineStatusSnapshotCreating        = "creating_snapshot"
	VirtualMachineStatusSnapshotRestoring       = "restoring_snapshot"
	VirtualMachineStatusSnapshotRestoringFailed = "restore_snapshot_failed"
)

type VirtualMachineResponse struct {
	Status         ResponseStatus `json:"status"`
	Message        string         `json:"message,omitempty"`
	VirtualMachine VirtualMachine `json:"virtual_machine"`
}

type VirtualMachinesResponse struct {
	Status          ResponseStatus   `json:"status"`
	Message         string           `json:"message,omitempty"`
	VirtualMachines []VirtualMachine `json:"virtual_machines"`
}

type NewVirtualMachineResponse struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message,omitempty"`
	Id      int            `json:"id"`
}

type CancelDatesResponse struct {
	Status      ResponseStatus `json:"status"`
	Message     string         `json:"message,omitempty"`
	CancelDates []time.Time    `json:"dates"`
}

func (service *VirtualMachineService) List() (*[]VirtualMachine, error) {
	var response VirtualMachinesResponse

	_, err := service.client.Get(virtualMachinesBasePath, &response)

	if err != nil {
		return nil, err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	machines := response.VirtualMachines

	for i := range machines {
		machines[i].client = service.client
	}

	return &machines, err
}

func (service *VirtualMachineService) Add(machine *VirtualMachine) (*VirtualMachine, error) {
	if err := machine.Validate(); err != nil {
		return NewVirtualMachine(service.client), err
	}

	payload := machine.Payload()

	var response NewVirtualMachineResponse

	_, err := service.client.Post(virtualMachinesBasePath, payload, &response)

	if err != nil {
		return NewVirtualMachine(service.client), err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	machine.Id = response.Id

	return machine, err
}

func (service *VirtualMachineService) AddFromSnapshot(machine *VirtualMachine, snapshot *Snapshot) (*VirtualMachine, error) {
	if err := machine.Validate(); err != nil {
		return NewVirtualMachine(service.client), err
	}

	payload := machine.Payload()

	payload.Add("snapshot", strconv.Itoa(snapshot.Id))

	var response NewVirtualMachineResponse

	_, err := service.client.Post(virtualMachinesBasePath, payload, &response)

	if err != nil {
		return NewVirtualMachine(service.client), err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	machine.Id = response.Id

	return machine, err
}

func (service *VirtualMachineService) View(machineId int) (*VirtualMachine, error) {
	var response VirtualMachineResponse

	_, err := service.client.Get(service.path(strconv.Itoa(machineId)), &response)

	if err != nil {
		return NewVirtualMachine(service.client), err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	machine := response.VirtualMachine

	machine.client = service.client

	return &machine, err
}

func (service *VirtualMachineService) Edit(machine *VirtualMachine) (*VirtualMachine, error) {
	if err := machine.Validate(); err != nil {
		return machine, err
	}

	payload := machine.GetChanges()

	var response StatusResponse

	_, err := service.client.Post(service.path(strconv.Itoa(machine.Id)), payload, &response)

	if err != nil {
		return machine, err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	return machine, err
}

func (service *VirtualMachineService) Cancel(machine *VirtualMachine, date *time.Time) (*VirtualMachine, error) {
	payload := &url.Values{
		"date": {date.Format("2006-01-02")},
	}

	var response StatusResponse

	_, err := service.client.Post(service.path(strconv.Itoa(machine.Id)+"/cancel"), payload, &response)

	if err != nil {
		return machine, err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	if err == nil {
		machine.Cancelled = date
	}

	return machine, err
}

func (service *VirtualMachineService) UndoCancellation(machine *VirtualMachine) error {
	payload := &url.Values{
		"date": {"0"},
	}

	var response StatusResponse

	_, err := service.client.Post(service.path(strconv.Itoa(machine.Id)+"/cancel"), payload, &response)

	if err != nil {
		return err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	return err
}

func (service *VirtualMachineService) GetCancelDates(machine *VirtualMachine) (*[]time.Time, error) {
	var response CancelDatesResponse

	_, err := service.client.Get(service.path(strconv.Itoa(machine.Id)+"/cancel"), &response)

	if err != nil {
		return nil, err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	return &response.CancelDates, err
}

func (service *VirtualMachineService) RunTask(task string, machine *VirtualMachine) (*VirtualMachine, error) {
	if !isValidTask(task) {
		return machine, NewInvalidTaskError(task)
	}

	var response StatusResponse

	_, err := service.client.Get(service.path(strconv.Itoa(machine.Id)+"/"+task), &response)

	if err != nil {
		return machine, err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	return machine, err
}

func (service *VirtualMachineService) Reinstall(machine *VirtualMachine) (*VirtualMachine, error) {
	payload := &url.Values{
		"template":          {strconv.Itoa(machine.Template.Id)},
		"reinstall":         {"true"},
		"confirm_reinstall": {"true"},
	}

	var response StatusResponse

	_, err := service.client.Post(service.path(strconv.Itoa(machine.Id)), payload, &response)

	if err != nil {
		return machine, err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	return machine, err
}

func (service *VirtualMachineService) CreateSnapshot(machine *VirtualMachine, name string, online bool, overwrite bool) (*Snapshot, error) {
	// TODO: Endpoint does not return Snapshot ID after creation

	payload := &url.Values{
		"name": {name},
	}

	if online {
		payload.Add("online", "true")
	}

	if overwrite {
		payload.Add("overwrite", "true")
	}

	var response StatusResponse

	_, err := service.client.Post(service.path(strconv.Itoa(machine.Id)+"/create_snapshot"), payload, &response)

	if err != nil {
		return NewSnapshot(service.client), err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	snapshot := NewSnapshot(service.client)

	snapshot.Name = name

	return snapshot, err
}

func (service *VirtualMachineService) RestoreSnapshot(machine *VirtualMachine, snapshot *Snapshot) (*VirtualMachine, error) {
	payload := &url.Values{
		"snapshot": {strconv.Itoa(snapshot.Id)},
	}

	var response StatusResponse

	_, err := service.client.Post(service.path(strconv.Itoa(machine.Id)+"/restore_snapshot"), payload, &response)

	if err != nil {
		return machine, err
	}

	if response.Status == ResponseError {
		err = NewApiError(response.Message)
	}

	return machine, err
}

func (service *VirtualMachineService) path(path string) string {
	return fmt.Sprintf("%s/%s", virtualMachinesBasePath, path)
}

func (machine *VirtualMachine) Create() error {
	_, err := machine.client.VirtualMachine.Add(machine)

	return err
}

func (machine *VirtualMachine) CreateFromSnapshot(snapshot *Snapshot) error {
	_, err := machine.client.VirtualMachine.AddFromSnapshot(machine, snapshot)

	return err
}

func (machine *VirtualMachine) Commit() error {
	if machine.Id == 0 {
		return NewVirtualMachineNotCreatedError()
	}

	_, err := machine.client.VirtualMachine.Edit(machine)

	return err
}

func (machine *VirtualMachine) Cancel(date *time.Time) error {
	if machine.Id == 0 {
		return NewVirtualMachineNotCreatedError()
	}

	cancelDate := machine.findValidCancelDate(date)

	if cancelDate == nil {
		return NewInvalidCancelDateError(date)
	}

	_, err := machine.client.VirtualMachine.Cancel(machine, date)

	return err
}

func (machine *VirtualMachine) findValidCancelDate(date *time.Time) *time.Time {
	cancelDates, _ := machine.GetCancelDates()

	for _, cancelDate := range *cancelDates {
		if date.Equal(cancelDate) {
			return date
		}
	}

	return nil
}

func (machine *VirtualMachine) UndoCancellation() error {
	if machine.Id == 0 {
		return NewVirtualMachineNotCreatedError()
	}

	if machine.Cancelled.IsZero() || machine.Cancelled.Before(time.Now()) {
		return NewVirtualMachineNotCancelledError()
	}

	err := machine.client.VirtualMachine.UndoCancellation(machine)

	return err
}

func (machine *VirtualMachine) GetCancelDates() (*[]time.Time, error) {
	if machine.Id == 0 {
		return nil, NewVirtualMachineNotCreatedError()
	}

	cancelDates, err := machine.client.VirtualMachine.GetCancelDates(machine)

	return cancelDates, err
}

func (machine *VirtualMachine) Start() error {
	if machine.Id == 0 {
		return NewVirtualMachineNotCreatedError()
	}

	_, err := machine.client.VirtualMachine.RunTask(TaskStart, machine)

	return err
}

func (machine *VirtualMachine) Stop() error {
	if machine.Id == 0 {
		return NewVirtualMachineNotCreatedError()
	}

	_, err := machine.client.VirtualMachine.RunTask(TaskStop, machine)

	return err
}

func (machine *VirtualMachine) Restart() error {
	if machine.Id == 0 {
		return NewVirtualMachineNotCreatedError()
	}

	_, err := machine.client.VirtualMachine.RunTask(TaskRestart, machine)

	return err
}

func (machine *VirtualMachine) PowerOff() error {
	if machine.Id == 0 {
		return NewVirtualMachineNotCreatedError()
	}

	_, err := machine.client.VirtualMachine.RunTask(TaskPowerOff, machine)

	return err
}

func (machine *VirtualMachine) Rescue() error {
	if machine.Id == 0 {
		return NewVirtualMachineNotCreatedError()
	}

	_, err := machine.client.VirtualMachine.RunTask(TaskRescue, machine)

	return err
}

func (machine *VirtualMachine) Reinstall() error {
	if machine.Id == 0 {
		return NewVirtualMachineNotCreatedError()
	}

	_, err := machine.client.VirtualMachine.Reinstall(machine)

	return err
}

func (machine *VirtualMachine) CreateSnapshot(name string, online bool, overwrite bool) (*Snapshot, error) {
	if machine.Id == 0 {
		return nil, NewVirtualMachineNotCreatedError()
	}

	return machine.client.VirtualMachine.CreateSnapshot(machine, name, online, overwrite)
}

func (machine *VirtualMachine) RestoreSnapshot(snapshot *Snapshot) error {
	if machine.Id == 0 {
		return NewVirtualMachineNotCreatedError()
	}

	if snapshot.Id == 0 {
		return NewSnapshotNotCreatedError()
	}

	_, err := machine.client.VirtualMachine.RestoreSnapshot(machine, snapshot)

	return err
}

func (machine *VirtualMachine) Payload() *url.Values {
	return &url.Values{
		"name":         {machine.Name},
		"dns_name":     {machine.Network[0].DnsName},
		"ram":          {strconv.Itoa(machine.Ram)},
		"storage":      {strconv.Itoa(machine.Storage.Size)},
		"storage_type": {string(machine.Storage.Type)},
		"template":     {strconv.Itoa(machine.Template.Id)},
		"site":         {strconv.Itoa(machine.Site.Id)},
		"cpu_count":    {strconv.Itoa(machine.Cpu.Cores)},
		"cpu_cap":      {strconv.Itoa(machine.Cpu.Cap)},
	}
}

func (machine *VirtualMachine) Validate() error {
	// TODO: Implement validation for all the fields, including by checking the permitted values for ram, storage, template, site and cpu
	return nil
}

func (machine *VirtualMachine) Refresh() error {
	update, err := machine.client.VirtualMachine.View(machine.Id)

	if err != nil {
		return err
	}

	machine.Name = update.Name
	machine.Cpu = update.Cpu
	machine.Ram = update.Ram
	machine.Storage = update.Storage
	machine.Site = update.Site
	machine.Template = update.Template
	machine.Network = update.Network
	machine.Admin = update.Admin
	machine.Managed = update.Managed
	machine.Locked = update.Locked
	machine.Status = update.Status
	machine.Created = update.Created
	machine.Cancelled = update.Cancelled

	return nil
}

func (machine *VirtualMachine) GetChanges() *url.Values {
	changes := &url.Values{}

	if machine.original.Name != "" {
		changes.Add("name", machine.Name)

		machine.original.Name = ""
	}

	if machine.original.Ram != 0 {
		changes.Add("ram", strconv.Itoa(machine.Ram))

		machine.original.Ram = 0
	}

	if machine.original.Storage != 0 {
		changes.Add("storage", strconv.Itoa(machine.Storage.Size))

		if machine.Storage.Size < machine.original.Storage {
			// TODO: Check if reinstall has to be set too
			changes.Add("confirm_reinstall", "true")
		}

		machine.original.Storage = 0
	}

	if machine.original.Cpu.Cores != 0 {
		changes.Add("cpu_count", strconv.Itoa(machine.Cpu.Cores))

		machine.original.Cpu.Cores = 0
	}

	if machine.original.Cpu.Cap != 0 {
		changes.Add("cpu_cap", strconv.Itoa(machine.Cpu.Cap))

		machine.original.Cpu.Cap = 0
	}

	return changes
}

func (machine *VirtualMachine) SetName(name string) error {
	if machine.Id == 0 {
		return NewVirtualMachineNotCreatedError()
	}

	if machine.original.Name == "" {
		machine.original.Name = machine.Name
	}

	machine.Name = name

	return nil
}

func (machine *VirtualMachine) SetRam(size int) error {
	if machine.Id == 0 {
		return NewVirtualMachineNotCreatedError()
	}

	if machine.original.Ram == 0 {
		machine.original.Ram = machine.Ram
	}

	machine.Ram = size

	return nil
}

func (machine *VirtualMachine) SetCpuCores(cores int) error {
	if machine.Id == 0 {
		return NewVirtualMachineNotCreatedError()
	}

	if machine.original.Cpu.Cores == 0 {
		machine.original.Cpu.Cores = machine.Cpu.Cores
	}

	machine.Cpu.Cores = cores

	return nil
}

func (machine *VirtualMachine) SetCpuCap(cap int) error {
	if machine.Id == 0 {
		return NewVirtualMachineNotCreatedError()
	}

	if machine.original.Cpu.Cap == 0 {
		machine.original.Cpu.Cap = machine.Cpu.Cap
	}

	machine.Cpu.Cap = cap

	return nil
}

func NewVirtualMachine(client *Client) *VirtualMachine {
	return &VirtualMachine{client: client}
}

func isValidTask(task string) bool {
	switch task {
	case
		TaskStart,
		TaskStop,
		TaskRestart,
		TaskRescue,
		TaskPowerOff:
		return true
	}

	return false
}
