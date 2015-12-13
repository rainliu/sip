package sip

type ServerTransaction interface {
	Transaction

	SendResponse(response Response) error
}
