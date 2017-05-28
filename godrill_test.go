package godrill

const (
  apiKey = ""
)

func TestGodrill(t *testing.T) {
  Key = apiKey

  email, err := NewTemplateEmail("test")
  if err != nil {
    t.Fatal(err)
  }

  email.SetSubject("Testing Godrill")
  email.SetFrom("godrilltester@godrill.com", "GodrillBot3000")
}
