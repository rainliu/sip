package sip

type ResponseEvent struct {
	transaction ClientTransaction
	response    Response
}

func NewResponseEvent(clientTransaction ClientTransaction, response Response) *ResponseEvent {
	return &ResponseEvent{
		transaction: clientTransaction,
		response:    response,
	}
}

func (this *ResponseEvent) GetClientTransaction() ClientTransaction {
	return this.transaction
}

func (this *ResponseEvent) GetResponse() Response {
	return this.response
}
