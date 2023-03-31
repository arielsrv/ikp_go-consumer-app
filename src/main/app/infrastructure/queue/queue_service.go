package queue

import "context"

type Service interface {
	Receive(ctx context.Context) ([]MessageDTO, error)
	Delete(ctx context.Context, receiptHandle string) error
	Count(ctx context.Context) (*int, error)
}

type MessageDTO struct {
	Body          string
	ReceiptHandle string
}

func (m *MessageDTO) String() string {
	return m.Body
}
