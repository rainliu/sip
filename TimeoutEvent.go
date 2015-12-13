package sip

type TimeoutEvent struct {
	transaction Transaction
	timeout     Timeout
}

func NewTimeoutEvent(transaction Transaction, timeout Timeout) *TimeoutEvent {
	return &TimeoutEvent{
		transaction: transaction,
		timeout:     timeout,
	}
}

func (this *TimeoutEvent) GetTransaction() Transaction {
	return this.transaction
}

func (this *TimeoutEvent) IsServerTransaction() bool {
	_, ok := this.transaction.(ServerTransaction)
	return ok
}

func (this *TimeoutEvent) GetTimeout() Timeout {
	return this.timeout
}
