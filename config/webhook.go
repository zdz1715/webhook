package config

const (
	WebhookDefaultMethod = "POST"
)

type KV struct {
	Name  string
	Value string
}

type Webhook struct {
	URL         string                 `yaml:"url"`
	Method      string                 `yaml:"method"`
	ContentType string                 `yaml:"contentType"`
	Header      []KV                   `yaml:"header"`
	Body        map[string]interface{} `yaml:"body"`
	Query       []KV                   `yaml:"query"`
	Vars        map[string]WebhookVar  `yaml:"varsFrom"`
	Client      *Client                `yaml:"client"`
}

const (
	WebhookVarFromQuery  = "Query"
	WebhookVarFromBody   = "Body"
	WebhookVarFromHeader = "Header"
	WebhookVarFromValue  = "Value"
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
