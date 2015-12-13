package sip

type ClientTransaction interface {
	Transaction

	SendRequest() error
	CreateCancel() (Request, error)
	CreateAck() (Request, error)
}
