package sip

type ClientTransaction interface {
	Transaction

	SendRequest() error
	CreateCancel() (Request, error)
	CreateAck() (Request, error)
}

type clientTransaction struct {
	transaction
}

func newClientTransaction(request Request) *clientTransaction {
	return &clientTransaction{
		transaction: transaction{
			request: request,
			quit:    make(chan bool),
		},
	}
}

func (this *clientTransaction) SendRequest() error {
	return nil
}

func (this *clientTransaction) CreateCancel() (Request, error) {
	return nil, nil
}

func (this *clientTransaction) CreateAck() (Request, error) {
	return nil, nil
}
