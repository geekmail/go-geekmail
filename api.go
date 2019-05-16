package geekmail

const (
	defaultBaseURL = "https://geekmail.app/api/1.0/"
	apiDraftCreate = "draft/create"
)

// Returned by the API
type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// TemplateMap is used to store variables used in the template
type TemplateMap map[string]string

// GitHubAuth contains the repository address and the secret.
type GitHubAuth struct {
	Repository string `json:"repository" yaml:"repository"` // e.g. github.com/geekmail/geekmail-sample
	Secret     string `json:"secret" yaml:"secret"`         // the secret found in geekmail.yaml
}

// payload for the "draft/create" endpoint
type draftCreate struct {
	// If fetching the template from GitHub
	GitHubAuth   GitHubAuth `json:"githubauth,omitempty"`
	TemplatePath string     `json:"template_path,omitempty"`
	// If the template is provided inline
	Template string `json:"template,omitempty"`

	Vars   TemplateMap `json:"vars,omitempty"`
	Labels []string    `json:"labels,omitempty"`
}
