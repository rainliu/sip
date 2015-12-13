package sip

type Transaction interface {
	GetDialog() Dialog
	GetState() TransactionState
	GetRetransmitTimer() int
	SetRetransmitTimer(retransmitTimer int)
	GetBranchId() string
	GetRequest() Request
	Close()
}

type TransactionState int

const (
	TRANSACTIONSTATE_CALLING    TransactionState = iota //0
	TRANSACTIONSTATE_TRYING                             //1
	TRANSACTIONSTATE_PROCEEDING                         //2
	TRANSACTIONSTATE_COMPLETED                          //3
	TRANSACTIONSTATE_CONFIRMED                          //4
	TRANSACTIONSTATE_TERMINATED                         //5
)
