package sip

type ServerTransaction interface {
	Transaction

	SendResponse(response Response) error
}

type serverTransaction struct {
	transaction

	response Response
}

func newServerTransaction() *serverTransaction {
	return &serverTransaction{}
}

func (this *serverTransaction) SendResponse(response Response) error {
	return nil
}
