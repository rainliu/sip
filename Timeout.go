package sip

const (
	TIMEOUT_RETRANSMIT  = iota //0
	TIMEOUT_TRANSACTION        //1
)

type Timeout struct {
	timeout int
}

func NewTimeout(timeout int) *Timeout {
	return &Timeout{timeout: timeout}
}

func (this *Timeout) GetValue() int {
	return this.timeout
}

func (this *Timeout) SetValue(timeout int) {
	this.timeout = timeout
}

func (this *Timeout) String() string {
	var text string
	switch this.timeout {
	case TIMEOUT_RETRANSMIT:
		text = "Retransmission Timeout"
	case TIMEOUT_TRANSACTION:
		text = "Transaction Timeout"
	default:
		text = "Error while printing Timeout"
	}
	return text
}
