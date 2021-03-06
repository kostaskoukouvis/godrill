package godrill

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

const apiURL = "https://mandrillapp.com/api/1.0/"

var (
	ErrNotInitialized     = errors.New("No Mandrill API key exists.")
	ErrNoTemplateName     = errors.New("Tried to add template content without first setting the template.")
	ErrOddNumberArguments = errors.New("The number of arguments given is not even.")
	ErrKeyNotString       = errors.New("The first argument of a key-value pair must be a string.")
)

// The API key, used globally
var Key string

type EmailSendResponse []EmailSendResponseItem

type EmailSendResponseItem struct {
	Email        string `json:"email"`
	Status       string `json:"status"`
	RejectReason string `json:"reject_reason"`
	ID           string `json:"_id"`
}

type EmailSendErrorResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

// request makes a http request to the mandrill API
func request(method, path string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", apiURL, path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, err
}
