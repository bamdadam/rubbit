package request

type PublishRequest struct {
	Topic        string `json:"topic"`
	Message      string `json:"message"`
	PublishDelay string `json:"publish_delay"`
	IsDelayed    bool   `json:"delayed"`
}

type GetMessageRequest struct {
	Topic string `json:"topic"`
}
