package interfaces

type Publisher interface {
	PublishMessage(eventType string, payload interface{}, metadata map[string]string) error
	CloseMessage() error
}
