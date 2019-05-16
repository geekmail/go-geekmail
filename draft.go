package geekmail

import "context"

type DraftService service

type DraftCreate struct {
	TemplatePath string      `json:"template_path,omitempty"`
	Template     string      `json:"template,omitempty"`
	Vars         TemplateMap `json:"vars,omitempty"`
	Labels       []string    `json:"labels,omitempty"`
}

func (d *DraftCreate) payload(conf *Conf) draftCreate {
	return draftCreate{
		// If fetching the template from GitHub. Copy GitHub auth data from conf
		GitHubAuth:   conf.GitHubAuth,
		TemplatePath: d.TemplatePath,

		// If the template is provided inline
		Template: d.Template,

		Vars:   d.Vars,
		Labels: d.Labels,
	}
}

func (s *DraftService) Create(ctx context.Context, data *DraftCreate) (*APIResponse, error) {
	req, err := s.client.NewRequest("POST", apiDraftCreate, data.payload(s.client.conf))
	if err != nil {
		return nil, err
	}

	response := new(APIResponse)
	_, err = s.client.Do(ctx, req, response)

	return response, err
}
