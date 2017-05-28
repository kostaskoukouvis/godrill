package godrill

import (
	"fmt"
	"testing"
)

const (
	apiKey = "xxx"
)

func TestGodrill(t *testing.T) {
	Key = apiKey

	email, err := NewTemplateEmail("testingtemplate")
	if err != nil {
		t.Fatal(err)
	}

	email.SetSubject("Testing Godrill")
	email.SetFrom("godrilltester@godrill.com", "GodrillBot3000")
	err = email.SetRecipient("masta13killah@yahoo.gr", "Masta Killah", "name")
	if err != nil {
		t.Fatal(err)
	}
	err = email.SetTemplateContent("content", "hello there", "footer", "Godrill production ltd.")
	if err != nil {
		t.Fatal(err)
	}
	err = email.SetGlobalMergeVars("name", "Slave!")
	if err != nil {
		t.Fatal(err)
	}
	res, errRes, err := email.Send()
	if err != nil {
		t.Fatal(err)
	}
	if errRes != nil {
		fmt.Println(errRes)
		t.Fail()
	}
	fmt.Println(res)
}
