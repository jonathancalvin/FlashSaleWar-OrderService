package messaging

type OrderConsumer struct{}

func (c *OrderConsumer) Handle(message []byte) error {
	// map event → use case
	return nil
}