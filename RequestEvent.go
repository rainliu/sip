package sip

type RequestEvent struct {
	transaction ServerTransaction
	request     Request
}

func NewRequestEvent(serverTransaction ServerTransaction, request Request) *RequestEvent {
	return &RequestEvent{
		transaction: serverTransaction,
		request:     request,
	}
}

func (this *RequestEvent) GetServerTransaction() ServerTransaction {
	return this.transaction
}

func (this *RequestEvent) GetRequest() Request {
	return this.request
}
