package godrill

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
)

type TemplateEmail struct {
	Key             string             `json:"key"`
	TemplateName    string             `json:"template_name"`
	TemplateContent []*TemplateContent `json:"template_content"`
	Message         *Message           `json:"message"`
}

type TemplateContent struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Message struct {
	Subject         string      `json:"subject"`
	To              []*To       `json:"to"`
	FromName        string      `json:"from_name"`
	FromEmail       string      `json:"from_email"`
	MergeVars       []*MergeVar `json:"merge_vars"`
	GlobalMergeVars []*Var      `json:"global_merge_vars"`
}

type To struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

type MergeVar struct {
	Recipient string `json:"rcpt"`
	Vars      []*Var `json:"vars"`
}

type Var struct {
	Name    string      `json:"name"`
	Content interface{} `json:"content"`
}

func NewTemplateEmail(name string) (*TemplateEmail, error) {
	if Key == "" {
		return nil, ErrNotInitialized
	}
	message := &Message{}
	email := &TemplateEmail{Key: Key, TemplateName: name, Message: message}
	return email, nil
}

func (e *TemplateEmail) SetTemplateContent(content ...string) error {
	if e.TemplateName == "" {
		return ErrNoTemplateName
	}
	contentArr, err := formatTemplateContent(content)
	if err != nil {
		return err
	}

	e.TemplateContent = contentArr
	return nil
}

func (e *TemplateEmail) SetSubject(subject string) {
	e.Message.Subject = subject
}

func (e *TemplateEmail) SetFrom(email, name string) {
	e.Message.FromEmail = email
	e.Message.FromName = name
}

func (e *TemplateEmail) SetRecipient(email, name string, args ...interface{}) error {
	to := &To{
		Email: email,
		Name:  name,
		Type:  "to",
	}
	e.Message.To = append(e.Message.To, to)
	if len(args) > 1 {
		vars, err := formatVar(args)
		if err != nil {
			return err
		}
		mergeVar := &MergeVar{
			Recipient: email,
			Vars:      vars,
		}
		e.Message.MergeVars = append(e.Message.MergeVars, mergeVar)
	}
	return nil
}

func (e *TemplateEmail) SetCC(email, name string) {
	to := &To{
		Email: email,
		Name:  name,
		Type:  "cc",
	}
	e.Message.To = append(e.Message.To, to)
}

func (e *TemplateEmail) SetBCC(email, name string) {
	to := &To{
		Email: email,
		Name:  name,
		Type:  "bcc",
	}
	e.Message.To = append(e.Message.To, to)
}

func (e *TemplateEmail) SetGlobalMergeVars(value ...interface{}) error {
	vars, err := formatVar(value)
	if err != nil {
		return err
	}

	e.Message.GlobalMergeVars = vars
	return nil
}

func (e *TemplateEmail) Send() (*EmailSendResponse, *EmailSendErrorResponse, error) {
	sendTemplatePath := "messages/send-template.json"
	emailBytes, err := json.Marshal(e)
	if err != nil {
		return nil, nil, err
	}
	emailBuf := bytes.NewBuffer(emailBytes)

	emailRes, err := request("POST", sendTemplatePath, emailBuf)
	if err != nil {
		return nil, nil, err
	}
	defer emailRes.Body.Close()

	decoder := json.NewDecoder(emailRes.Body)

	if emailRes.StatusCode != http.StatusOK {
		var errRes EmailSendErrorResponse
		err = decoder.Decode(&errRes)
		if err != nil {
			return nil, nil, err
		}

		return nil, &errRes, nil
	}

	var er EmailSendResponse
	err = decoder.Decode(&er)
	if err != nil {
		return nil, nil, err
	}

	return &er, nil, nil
}

// formatTemplateContent creates an array of TemplateContent items
// from then given string array. It assumes that the first value
// of the pair will be the name while the second the content of
// the templates editable region.
func formatTemplateContent(arr []string) ([]*TemplateContent, error) {
	if len(arr)%2 != 0 {
		return nil, ErrOddNumberArguments
	}
	result := []*TemplateContent{}
	// work the array in pairs.
	for i := 0; i < len(arr); i += 2 {
		pair := &TemplateContent{
			Name:    arr[i],
			Content: arr[i+1],
		}
		result = append(result, pair)
	}
	return result, nil
}

// formatVar creates an array of Var structs from the given
// interface{} array. It assumes that the first value of the
// pair will be the name while the second the content
func formatVar(arr []interface{}) ([]*Var, error) {
	if len(arr)%2 != 0 {
		return nil, ErrOddNumberArguments
	}
	result := []*Var{}
	// work the array in pairs.
	for i := 0; i < len(arr); i += 2 {
		// unless the first value is a string, dump the pair.
		v := reflect.TypeOf(arr[i]).Kind()
		if v != reflect.String {
			return nil, ErrKeyNotString
		}
		pair := &Var{
			Name:    arr[i].(string),
			Content: arr[i+1],
		}
		result = append(result, pair)
	}
	return result, nil
}
