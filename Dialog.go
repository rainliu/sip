package sip

type Dialog interface {
	GetLocalParty() string
	GetRemoteParty() string
	GetRemoteTarget() string
	GetDialogId() string
	GetCallId() string
	GetLocalSequenceNumber() int
	GetRemoteSequenceNumber() int
	GetRouteSet() []string
	IsSecure() bool
	IsServer() bool
	IncrementLocalSequenceNumber()
	CreateRequest(method string) (Request, error)
	SendRequest(ct ClientTransaction) error
	SendAck(ack Request) error
	GetState() DialogState
	Close()
	GetFirstTransaction() Transaction
	GetLocalTag() string
	GetRemoteTag() string
	SetApplicationData(applicationData interface{})
	GetApplicationData() interface{}
}

type DialogState int

const (
	DIALOGSTATE_EARLY      DialogState = iota //0
	DIALOGSTATE_CONFIRMED                     //1
	DIALOGSTATE_COMPLETED                     //2
	DIALOGSTATE_TERMINATED                    //3
)
