package geekmail

import (
	"context"
	"net/http"
	"testing"

	"github.com/akfaew/test"
)

func TestDraftCreate(t *testing.T) {
	c := NewClient(&http.Client{}, &TestConf)

	template := test.InputFixture(t, "subject.template")
	resp, err := c.Draft.Create(context.Background(), &DraftCreate{
		Template: string(template),

		Vars: TemplateMap{
			"To":      "John Doe <john@example.com>",
			"Subject": t.Name(),
		},
		Labels: []string{"GeekMail"},
	})

	test.NoError(t, err)
	resp.Message = "const"
	test.Fixture(t, resp)
}

func TestDraftCreateGitHub(t *testing.T) {
	c := NewClient(&http.Client{}, &TestConf)

	resp, err := c.Draft.Create(context.Background(), &DraftCreate{
		TemplatePath: "templates/subject.template",

		Vars: TemplateMap{
			"To":      "John Doe <john@example.com>",
			"Subject": t.Name(),
		},
		Labels: []string{"GeekMail"},
	})

	test.NoError(t, err)
	resp.Message = "const"
	test.Fixture(t, resp)
}
