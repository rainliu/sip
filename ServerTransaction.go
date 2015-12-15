package sip

type ServerTransaction interface {
	Transaction

	SendResponse(Response) error
}

type serverTransaction struct {
	transaction
}

func newServerTransaction() *serverTransaction {
	return &serverTransaction{
		transaction: transaction{
			quit: make(chan bool),
		},
	}
}

func (this *serverTransaction) SendResponse(resp Response) error {
	return nil
}
