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

///////////////////////////////////////////////////////////////
type transaction struct {
	dialog           Dialog
	transactionState TransactionState
	retransmitTimer  int
	branchId         string
	request          Request
	quit             chan bool
}

func (this *transaction) GetDialog() Dialog {
	return this.dialog
}
func (this *transaction) SetDialog(dialog Dialog) {
	this.dialog = dialog
}
func (this *transaction) GetState() TransactionState {
	return this.transactionState
}
func (this *transaction) SetState(transactionState TransactionState) {
	this.transactionState = transactionState
}
func (this *transaction) GetRetransmitTimer() int {
	return this.retransmitTimer
}
func (this *transaction) SetRetransmitTimer(retransmitTimer int) {
	this.retransmitTimer = retransmitTimer
}
func (this *transaction) GetBranchId() string {
	return this.branchId
}
func (this *transaction) SetBranchId(branchId string) {
	this.branchId = branchId
}
func (this *transaction) GetRequest() Request {
	return this.request
}
func (this *transaction) Close() {
	close(this.quit)
}
