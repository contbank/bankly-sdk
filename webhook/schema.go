package webhook

type ConfigItem struct {
	Name      string `json:"name"`
	Context   string `json:"context"`
	EventName string `json:"eventName"`
	Uri       string `json:"uri"`
	PublicKey string `json:"publicKey"`
}

type ConfigEntity struct {
	Id string `json:"id"`
	ConfigItem
}

type RegisterWebhookRequest struct {
	ConfigItem
	PrivateKey string `json:"privateKey"`
}

type SchemaLink struct {
	Url    string `json:"url"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type RegisterWebhookResponse struct {
	Data  ConfigEntity `json:"data"`
	Links []SchemaLink `json:"links"`
}
