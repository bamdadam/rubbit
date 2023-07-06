package request

type PublishRequest struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}
