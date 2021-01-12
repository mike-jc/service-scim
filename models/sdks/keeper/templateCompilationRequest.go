package modelsSdkKeeper

type TemplateCompilationRequest struct {
	TemplateName string                 `json:"templateName"`
	Compile      bool                   `json:"compile"`
	User         map[string]interface{} `json:"user,omitempty"`
	Group        map[string]interface{} `json:"group,omitempty"`
}
