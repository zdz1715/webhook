package config

const (
	WebhookDefaultMethod = "POST"
)

type KV struct {
	Name  string
	Value string
}

type Webhook struct {
	URL         string                `yaml:"url"`
	Method      string                `yaml:"method"`
	ContentType string                `yaml:"contentType"`
	Header      []KV                  `yaml:"header"`
	Body        WebhookBody           `yaml:"body"`
	Query       []KV                  `yaml:"query"`
	Vars        map[string]WebhookVar `yaml:"varsFrom"`
	Client      *Client               `yaml:"client"`
}

type WebhookBody struct {
	Form map[string]string `yaml:"form"`
	Json string            `json:"json"`
}

const (
	WebhookVarFromQuery  = "Query"
	WebhookVarFromBody   = "Body"
	WebhookVarFromHeader = "Header"
)

//var WebhookVarFrom map[string]struct{} = map[string]struct{}{
//	WebhookVarFromQuery:  struct{}{},
//	WebhookVarFromBody:   struct{}{},
//	WebhookVarFromHeader: struct{}{},
//	WebhookVarFromValue:  struct{}{},
//}

type WebhookVar struct {
	From  string `yaml:"from"`
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}
