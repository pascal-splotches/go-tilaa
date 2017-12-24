package go_tilaa

import (
	"fmt"
	"net/http"
	"time"
)

type ApiError struct {
	reason string
}

type ApiRequestError struct {
	reason string
}

type ClientError struct {
	reason string
}

type InvalidCredentialsError struct {
	credentials BasicAuth
}

type ResultsDecoderError struct {
	jsonDecoderError error
	response         *http.Response
}

type InvalidTaskError struct {
	task string
}

type InvalidCancelDateError struct {
	date *time.Time
}

type VirtualMachineNotCreatedError struct {
}

type VirtualMachineNotCancelledError struct {
}

type SnapshotNotCreatedError struct {
}

type MetadataNotCreatedError struct {
}

type SshKeyNotCreatedError struct {
}

var _ error = &ApiError{}
var _ error = &ApiRequestError{}
var _ error = &ClientError{}

var _ error = &InvalidCredentialsError{}
var _ error = &ResultsDecoderError{}

var _ error = &InvalidTaskError{}
var _ error = &InvalidCancelDateError{}

var _ error = &VirtualMachineNotCreatedError{}
var _ error = &VirtualMachineNotCancelledError{}
var _ error = &SnapshotNotCreatedError{}
var _ error = &MetadataNotCreatedError{}
var _ error = &SshKeyNotCreatedError{}

func NewApiError(reason string) *ApiError {
	return &ApiError{reason: reason}
}

func NewApiRequestError(reason string) *ApiRequestError {
	return &ApiRequestError{reason: reason}
}

func NewClientError(reason string) *ClientError {
	return &ClientError{reason: reason}
}

func NewInvalidCredentialsError(credentials BasicAuth) *InvalidCredentialsError {
	return &InvalidCredentialsError{credentials: credentials}
}

func NewResultsDecoderError(jsonDecoderError error, response *http.Response) *ResultsDecoderError {
	return &ResultsDecoderError{jsonDecoderError: jsonDecoderError, response: response}
}

func NewInvalidTaskError(task string) *InvalidTaskError {
	return &InvalidTaskError{task: task}
}

func NewInvalidCancelDateError(date *time.Time) *InvalidCancelDateError {
	return &InvalidCancelDateError{date: date}
}

func NewVirtualMachineNotCreatedError() *VirtualMachineNotCreatedError {
	return &VirtualMachineNotCreatedError{}
}

func NewVirtualMachineNotCancelledError() *VirtualMachineNotCancelledError {
	return &VirtualMachineNotCancelledError{}
}

func NewSnapshotNotCreatedError() *SnapshotNotCreatedError {
	return &SnapshotNotCreatedError{}
}

func NewMetadataNotCreatedError() *MetadataNotCreatedError {
	return &MetadataNotCreatedError{}
}

func NewSshKeyNotCreatedError() *SshKeyNotCreatedError {
	return &SshKeyNotCreatedError{}
}

func (error *ApiError) Error() string {
	return fmt.Sprintf("API Error: %s", error.reason)
}

func (error *ApiRequestError) Error() string {
	return fmt.Sprintf("API Request Error: %s", error.reason)
}

func (error *ClientError) Error() string {
	return fmt.Sprintf("Client Error: %s", error.reason)
}

func (error *InvalidCredentialsError) Error() string {
	return fmt.Sprintf("Invalid API Credentials: Username (%s) or Password (%s) invalid", error.credentials.UserName, error.credentials.Password)
}

func (error *ResultsDecoderError) Error() string {
	return fmt.Sprintf("Results Decoder Error: Unable to decode result (%s)", error.jsonDecoderError.Error())
}

func (error *InvalidTaskError) Error() string {
	return fmt.Sprintf("Invalid Task: %s", error.task)
}

func (error *InvalidCancelDateError) Error() string {
	return fmt.Sprintf("Invalid Cancel Date: %s", error.date.String())
}

func (error *VirtualMachineNotCreatedError) Error() string {
	return fmt.Sprintf("Virtual Machine has not been created yet.")
}

func (error *VirtualMachineNotCancelledError) Error() string {
	return fmt.Sprintf("Virtual Machine was not cancelled.")
}

func (error *SnapshotNotCreatedError) Error() string {
	return fmt.Sprintf("Snapshot has not been created yet.")
}

func (error *MetadataNotCreatedError) Error() string {
	return fmt.Sprintf("Metadata has not been created yet.")
}

func (error *SshKeyNotCreatedError) Error() string {
	return fmt.Sprintf("SshKey has not been created yet.")
}