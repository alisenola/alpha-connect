package messages

type PubSubMessage struct {
	Topic   string
	ID      uint64
	Message interface{}
}
