package webhook

type RegisterWebhookRequest struct {
	Name       string `json:"name"`
	EventName  string `json:"eventName"`
	Context    string `json:"context"`
	Uri        string `json:"uri"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

type SchemaLink struct {
	Url    string `json:"url"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type WebhookRecord struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Context   string `json:"context"`
	EventName string `json:"eventName"`
	Uri       string `json:"uri"`
	PublicKey string `json:"publicKey"`
}

type RegisterWebhookResponse struct {
	Data  WebhookRecord `json:"data"`
	Links []SchemaLink  `json:"links"`
}
